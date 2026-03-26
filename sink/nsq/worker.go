package nsq

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	phpSerialize "github.com/yvasiyarov/php_session_decoder/php_serialize"
)

const (
	actionWrite  = "write"
	actionUpdate = "update"
	actionDelete = "delete"
)

// NsqWorker 处理单张表的消息，将 OriginRow 格式化后发送到对应 topic
type NsqWorker struct {
	DbName     string
	TableName  string
	TopicNames []string
	isNewTable bool
	pool       *NsdPool
	Reader     chan *model.OriginRow
	stop       chan bool
}

func NewNsqWorker(dbName, tableName string, topicNames []string, isNewTable bool, pool *NsdPool) *NsqWorker {
	return &NsqWorker{
		DbName:     dbName,
		TableName:  tableName,
		TopicNames: topicNames,
		isNewTable: isNewTable,
		pool:       pool,
		Reader:     make(chan *model.OriginRow, 100),
		stop:       make(chan bool, 1),
	}
}

func (w *NsqWorker) Start(wg *sync.WaitGroup) {
	log.Printf("[NsqWorker] %s.%s started", w.DbName, w.TableName)
	wg.Add(1)
	defer func() {
		wg.Done()
		log.Printf("[NsqWorker] %s.%s stopped", w.DbName, w.TableName)
	}()
	for {
		select {
		case row := <-w.Reader:
			w.process(row)
		case <-w.stop:
			return
		}
	}
}

func (w *NsqWorker) Stop() {
	w.stop <- true
}

func (w *NsqWorker) process(row *model.OriginRow) {
	var action string
	var rowData *map[string]string

	switch row.Op {
	case model.OpType_Write:
		action = actionWrite
		rowData = row.After
	case model.OpType_Update:
		action = actionUpdate
		rowData = row.After
	case model.OpType_Delete:
		action = actionDelete
		rowData = row.Before
	default:
		// OpType_Update_Row 等不转发到 NSQ
		return
	}

	if rowData == nil {
		return
	}

	var msg string
	var err error
	if w.isNewTable {
		msg, err = w.buildJSON(action, row, rowData)
	} else {
		msg, err = w.buildPHP(action, row, rowData)
	}
	if err != nil {
		log.Printf("[NsqWorker] build msg error table=%s err=%v", w.TableName, err)
		return
	}

	data := []byte(msg)
	for _, topic := range w.TopicNames {
		go w.pool.Send(w.TableName, topic, data)
	}
}

// buildJSON 生成 JSON 格式消息（新表）
func (w *NsqWorker) buildJSON(action string, row *model.OriginRow, rowData *map[string]string) (string, error) {
	var tableData interface{}
	if action == actionUpdate && row.Before != nil {
		tableData = map[string]interface{}{
			"before": *row.Before,
			"after":  *rowData,
		}
	} else {
		tableData = *rowData
	}

	data := map[string]interface{}{
		"db":        w.DbName,
		"table":     w.TableName,
		"type":      action,
		"timestamp": row.Now / 1e6, // nanoseconds → milliseconds
		"data":      []interface{}{tableData},
	}
	b, err := json.Marshal(data)
	return string(b), err
}

// buildPHP 生成 PHP 序列化格式消息（旧表）
func (w *NsqWorker) buildPHP(action string, row *model.OriginRow, rowData *map[string]string) (string, error) {
	var tableData phpSerialize.PhpValue
	if action == actionUpdate && row.Before != nil {
		tableData = phpSerialize.PhpArray{
			"before": mapStringToPhpArray(*row.Before),
			"after":  mapStringToPhpArray(*rowData),
		}
	} else {
		tableData = mapStringToPhpArray(*rowData)
	}

	data := phpSerialize.PhpArray{
		"db":        phpSerialize.PhpValueString(w.DbName),
		"table":     phpSerialize.PhpValueString(w.TableName),
		"type":      phpSerialize.PhpValueString(action),
		"timestamp": phpSerialize.PhpValueUInt(uint64(row.Now / 1e6)),
		"data": phpSerialize.PhpSlice{
			tableData,
		},
	}
	return phpSerialize.Serialize(data)
}

func mapStringToPhpArray(m map[string]string) phpSerialize.PhpArray {
	arr := make(phpSerialize.PhpArray, len(m))
	for k, v := range m {
		// 去掉尾部空字节（MySQL binary 字段可能含有 \x00）
		arr[k] = phpSerialize.PhpValueString(strings.TrimRight(v, "\x00"))
	}
	return arr
}
