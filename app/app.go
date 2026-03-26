package app

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/complement"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/mq"
	"github.com/olachat/banban_server/cli_mysql_pipeline/sink"
)

type App struct {
	data         []*index.Index
	sinks        []sink.Sink
	wg           *sync.WaitGroup
	stopKafka    chan bool
	stopNsq      chan bool
	receiver     chan *model.OriginRow
	indexMap     map[string][]*index.Index
	sinkMap      map[string][]sink.Sink
	fullComplete bool
	processNum   int64
	cancel       chan bool
}

var currentApp *App

func NewApp() *App {
	if currentApp == nil {
		currentApp = &App{
			data:      make([]*index.Index, 0),
			sinks:     make([]sink.Sink, 0),
			wg:        &sync.WaitGroup{},
			stopKafka: make(chan bool, 0),
			stopNsq:   make(chan bool, 0),
			receiver:  make(chan *model.OriginRow, 0),
			indexMap:  make(map[string][]*index.Index),
			sinkMap:   make(map[string][]sink.Sink),
			cancel:    make(chan bool, 0),
		}
	}
	return currentApp
}

// Add 注册一个 ES Index
func (app *App) Add(item *index.Index) {
	app.data = append(app.data, item)
}

// RegisterSink 注册一个 Sink（如 NSQ Sink）
func (app *App) RegisterSink(s sink.Sink) {
	app.sinks = append(app.sinks, s)
	for _, table := range s.Tables() {
		app.sinkMap[table] = append(app.sinkMap[table], s)
	}
}

func (app *App) IsFullComplete() bool {
	return app.fullComplete
}

func (app *App) checkName(val string, indexNames map[string]bool) map[string]bool {
	names := strings.Split(val, ",")
	data := map[string]bool{}
	for _, name := range names {
		if _, ok := indexNames[name]; !ok {
			panic(fmt.Errorf("error args with start"))
		}
		data[name] = true
	}
	return data
}

func (app *App) Start(startWithFull string, dropIndexName string) {
	app.fullComplete = false
	indexNames := map[string]bool{}
	tableNames := map[string]bool{}

	for i := 0; i < len(app.data); i++ {
		app.data[i].Init(app.wg)
		for j := 0; j < len(app.data[i].Ops); j++ {
			name := fmt.Sprintf(
				"%s.%s",
				app.data[i].Ops[j].Db,
				app.data[i].Ops[j].Name,
			)
			if _, ok := app.indexMap[name]; !ok {
				app.indexMap[name] = make([]*index.Index, 0)
			}
			app.indexMap[name] = append(app.indexMap[name], app.data[i])
			tableNames[name] = true
		}
		indexNames[app.data[i].Name] = true
	}

	// 同时把 NSQ Sink 关心的表也加入 tableNames，让 Kafka consumer 不过滤这些表
	for table := range app.sinkMap {
		tableNames[table] = true
	}

	for i := 0; i < len(app.data); i++ {
		go app.data[i].Listen()
	}

	// 启动所有 Sink
	for _, s := range app.sinks {
		s.Start()
	}

	if len(startWithFull) > 0 {
		fullNames := app.checkName(startWithFull, indexNames)
		dropIndexNames := map[string]bool{}
		if len(dropIndexName) > 0 {
			dropIndexNames = app.checkName(dropIndexName, indexNames)
		}
		for i := 0; i < len(app.data); i++ {
			if _, ok := fullNames[app.data[i].Name]; ok {
				_, ok := dropIndexNames[app.data[i].Name]
				app.data[i].Full(ok)
			}
		}
		glg.Infof("full init complete")
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}
	app.fullComplete = true

	defer func() {
		glg.Infof("app exit")
		app.wg.Done()
	}()
	glg.Infof("begin to consume mq")
	glg.Infof("tableNames: %v", tableNames)

	go func() {
		defer app.StopIndexListen()
		mq.NewKafkaMq(app.stopKafka, app.receiver, tableNames)
	}()

	go mq.NewNsqMq(app.stopNsq, app.receiver, mq.DefaultConverter)

	for {
		select {
		case message := <-complement.Instance.Receiver:
			app.onMessage(message)
		case message := <-app.receiver:
			app.onMessage(message)
		case <-app.cancel:
			return
		}
	}
}

func (app *App) StopIndexListen() {
	for i := 0; i < len(app.data); i++ {
		app.wg.Add(1)
		go app.data[i].Stop()
	}
	app.cancel <- true
}

func (app *App) onMessage(message *model.OriginRow) {
	atomic.AddInt64(&app.processNum, 1)
	name := message.Table

	// 路由到 ES Index
	if idxs, ok := app.indexMap[name]; ok {
		for i := 0; i < len(idxs); i++ {
			idxs[i].Receiver <- message
		}
	}

	// 路由到 Sink（如 NSQ Sink）
	if sinks, ok := app.sinkMap[name]; ok {
		for _, s := range sinks {
			s.Send(message)
		}
	}
}

func (app *App) Stop() {
	app.wg.Add(1)
	glg.Infof("Stop stopNsq")
	app.stopNsq <- true
	glg.Infof("Stop stopKafka")
	if !conf.IsDev {
		app.stopKafka <- true
	}
	// 停止所有 Sink
	for _, s := range app.sinks {
		s.Stop()
	}
	app.cancel <- true
	app.wg.Wait()
	glg.Infof("Stop exit")
}

func (app *App) GetProcessNum() int64 {
	return atomic.SwapInt64(&app.processNum, 0)
}
