package nsq

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
)

// NsqSink 实现 sink.Sink 接口，将 OriginRow 消息按表路由到对应 NSQ Worker
type NsqSink struct {
	workers map[string]*NsqWorker // "db.table" → worker
	tables  []string
	pool    *NsdPool
	wg      sync.WaitGroup
}

func NewNsqSink(cfg conf.NsqPublisherConf) *NsqSink {
	pool := NewNsdPool(cfg.Pools, cfg.Topics, cfg.Debug)

	newTableSet := make(map[string]bool, len(cfg.NewTables))
	for _, t := range cfg.NewTables {
		newTableSet[t] = true
	}

	workers := make(map[string]*NsqWorker)
	tables := make([]string, 0, len(cfg.Tables))

	for tableName, topicNames := range cfg.Tables {
		isNewTable := newTableSet[tableName]
		// tableName 可能是 "db.table" 或纯 "table"
		dbName, tName := splitDbTable(tableName)
		key := fmt.Sprintf("%s.%s", dbName, tName)
		workers[key] = NewNsqWorker(dbName, tName, topicNames, isNewTable, pool)
		tables = append(tables, key)
	}

	return &NsqSink{
		workers: workers,
		tables:  tables,
		pool:    pool,
	}
}

func (s *NsqSink) Name() string { return "nsq" }

func (s *NsqSink) Tables() []string { return s.tables }

func (s *NsqSink) Send(row *model.OriginRow) {
	if w, ok := s.workers[row.Table]; ok {
		w.Reader <- row
	}
}

func (s *NsqSink) Start() {
	for _, w := range s.workers {
		go w.Start(&s.wg)
	}
	log.Printf("[NsqSink] started %d workers", len(s.workers))
}

func (s *NsqSink) Stop() {
	for _, w := range s.workers {
		w.Stop()
	}
	s.wg.Wait()
	s.pool.Stop()
	log.Println("[NsqSink] stopped")
}

// splitDbTable 将 "db.table" 拆分为 (db, table)，如果没有点则默认 db="xianshi"
func splitDbTable(name string) (string, string) {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "xianshi", name
}
