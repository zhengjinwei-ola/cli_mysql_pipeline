package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/app"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/instance/partystar/cities_index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/instance/partystar/room_new"
	"github.com/olachat/banban_server/cli_mysql_pipeline/instance/partystar/song_index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/instance/partystar/user_god_new"
	"github.com/olachat/banban_server/cli_mysql_pipeline/instance/partystar/user_new"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	_ "github.com/olachat/banban_server/cli_mysql_pipeline/mq"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	nsqsink "github.com/olachat/banban_server/cli_mysql_pipeline/sink/nsq"
	"github.com/olachat/banban_server/common/metrics"
	"github.com/olachat/banban_server/common/monitor"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	var startWithFull string
	var dropIndexName string
	var env string

	flag.StringVar(&env, "env", "partystar", "操作哪个配置文件数据")
	flag.StringVar(&startWithFull, "startWithFull", "", "哪些索引全量开始")
	flag.StringVar(&dropIndexName, "dropIndexName", "", "哪些索引先删除再重建")
	flag.Parse()

	glg.Info("start parse config")
	fmt.Println("env:", env)
	cfg := conf.ParseConfig(env)

	metrics.InitUpMetrics("bbsvr_cli_mysql_pipeline", "cmd")
	err := metrics.RegisterService(metrics.CmdServiceTag, metrics.CmdServiceName, metrics.DiscoverConfig{
		Type: "consul",
		Addr: []string{conf.Discovery},
		Path: "/banban",
	})
	if err != nil {
		panic(err)
	}

	glg.Info("start init db")
	model.InitDb()

	glg.Info("start upgrade mapping")
	tables.UpgradeOverseaUserMapping()

	glg.Info("start new app")
	control := app.NewApp()

	song := song_index.GetIndex()
	room := room_new.GetIndex()
	user := user_new.GetIndex()
	userGod := user_god_new.GetIndex()
	cities := cities_index.GetIndex()

	control.Add(song)
	control.Add(userGod)
	control.Add(user)
	control.Add(room)
	control.Add(cities)

	// 注册 NSQ Sink（配置了 nsq_publisher 才启用）
	if cfg.NsqPublisher != nil {
		control.RegisterSink(nsqsink.NewNsqSink(*cfg.NsqPublisher))
		glg.Info("NSQ Sink registered")
	}

	go control.Start(startWithFull, dropIndexName)

	monitor.PostMonitorMsg("banban_server.cli_mysql_pipeline 启动")

	setupBatchImport(room, user)

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	timerLog := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-timerLog.C:
			num := control.GetProcessNum()
			glg.Info("[main] ProcessNum:", num)
		case s := <-sign:
			glg.Infof("receive signal %d", s)
			if control.IsFullComplete() {
				control.Stop()
			}
			time.Sleep(time.Second * 1)
			monitor.PostMonitorMsg("banban_server.cli_mysql_pipeline 退出")
			err = metrics.DeRegisterService()
			if err != nil {
				g.Log().Errorf("DeRegisterService error:%s", err.Error())
			}
			return
		}
	}
}
