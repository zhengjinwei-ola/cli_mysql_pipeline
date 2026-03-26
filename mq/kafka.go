package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/olachat/banban_server/cli_mysql_pipeline/avro"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/tidwall/gjson"
)

// CanalJSON 定义DTS的数据结构
type CanalJSON struct {
	Data      []map[string]string `json:"data"`
	Database  string              `json:"database"`
	Table     string              `json:"table"`
	Es        int64               `json:"es"`        //操作在源库的执行时间，13位Unix时间戳，单位为毫秒
	Ts        int64               `json:"ts"`        //操作写入到目标库的时间，13位Unix时间戳，单位为毫秒
	Op        string              `json:"type"`      //操作的类型，比如DELETE、UPDATE、INSERT
	ID        int64               `json:"id"`        //操作的序列号
	IsDdl     bool                `json:"isDdl"`     //是否是DDL操作
	MysqlType map[string]string   `json:"mysqlType"` //字段的数据类型
	Old       []map[string]string `json:"old"`       //变更前的数据
	PkNames   []string            `json:"pkNames"`   //主键名称
	SQL       string              `json:"sql"`       //SQL语句
	SQLType   map[string]int      `json:"sqlType"`   //经转换处理后的字段类型
}

// ParseCanalJSON 解析[]byte
func ParseCanalJSON(data []byte) (*CanalJSON, error) {
	value := &CanalJSON{}
	err := json.Unmarshal(data, value)
	return value, err
}

const (
	CanalWrite  = "INSERT"
	CanalDelete = "DELETE"
	CanalUpdate = "UPDATE"
)

var opTypeMap map[string]model.OpType = map[string]model.OpType{
	CanalWrite:  model.OpType_Write,
	CanalUpdate: model.OpType_Update,
	CanalDelete: model.OpType_Delete,
}

func init() {
	compiler.LoggingEnabled = false
}

// 阿里的DTS莫名其妙，你没订阅的数据有时也在
func NewKafkaMq(stop chan bool, message chan *model.OriginRow, tableNames map[string]bool) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_2_0_0
	config.Consumer.Return.Errors = true
	// config.Net.SASL.Enable = true
	// config.Net.SASL.User = conf.MqUserName
	// config.Net.SASL.Password = conf.MqPassword
	config.Net.MaxOpenRequests = 100
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.MaxWaitTime = 300 * time.Millisecond

	group, err := sarama.NewConsumerGroup(
		[]string{conf.MqAddr},
		conf.MqGroup,
		config,
	)
	if err != nil {
		panic(err)
	}

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	consumer := NewConsumer(message, tableNames)
	topics := []string{conf.MqTopic}
	go func() {
		defer wg.Done()
		for {
			log.Println("group.Consume")
			err := group.Consume(ctx, topics, consumer)
			if err != nil {
				log.Println("group.break", err)
				break
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
		log.Println("group.end")
	}()
	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	<-stop
	log.Println("to close")
	if err = group.Close(); err != nil {
		log.Printf("Error closing client: %v\n", err)
	}
	log.Println("kafka closed")
	log.Println("wg.Wait")
	wg.Wait()

	log.Println("wg.Wait end")
}

func NewConsumer(message chan *model.OriginRow, tableNames map[string]bool) Consumer {
	record := avro.NewRecord()
	deser, err := compiler.CompileSchemaBytes([]byte(record.Schema()), []byte(record.Schema()))
	if err != nil {
		panic(err)
	}
	return Consumer{
		message:    message,
		tableNames: tableNames,
		ready:      make(chan bool),
		deser:      deser,
		record:     record,
	}
}

type Consumer struct {
	message    chan *model.OriginRow
	tableNames map[string]bool
	ready      chan bool
	deser      *vm.Program
	record     *avro.Record
}

func (c Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}
func (Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	findErrorTable := false
	for msg := range claim.Messages() {
		jsonMsg, err := ParseCanalJSON(msg.Value)
		if err != nil {
			log.Printf("the msg is not valid json. msg: %s", string(msg.Value))
			sess.MarkMessage(msg, "")
			msgJson := string(msg.Value)
			msgAfter := map[string]string{}
			msgBefore := map[string]string{}
			a := gjson.Get(msgJson, "data.0").String()
			b := gjson.Get(msgJson, "old.0").String()
			opType := gjson.Get(msgJson, "type").String()
			opTypeMap := map[string]model.OpType{
				"WRITE":  model.OpType_Write,
				"UPDATE": model.OpType_Update,
				"DELETE": model.OpType_Delete,
			}
			opMsg := opTypeMap[opType]
			json.Unmarshal([]byte(a), &msgAfter)
			json.Unmarshal([]byte(b), &msgBefore)
			name := fmt.Sprintf("%s.%s", gjson.Get(msgJson, "database").String(), gjson.Get(msgJson, "table").String())
			now := time.Now().UnixNano()
			message := &model.OriginRow{
				Table:  name,
				Op:     opMsg,
				Before: &msgBefore,
				After:  &msgAfter,
				Now:    now,
			}
			c.message <- message
			continue
		}
		log.Printf("jinweiDebug msg.Value: %s\n", string(msg.Value))
		name := fmt.Sprintf("%s.%s", jsonMsg.Database, jsonMsg.Table)
		if _, ok := c.tableNames[name]; !ok {
			findErrorTable = true
			sess.MarkMessage(msg, "")
			continue
		}

		op, ok := opTypeMap[jsonMsg.Op]
		if ok {
			before, after := make(map[string]string), make(map[string]string)
			if len(jsonMsg.Old) > 0 {
				before = jsonMsg.Old[0]
			}
			if len(jsonMsg.Data) > 0 {
				after = jsonMsg.Data[0]
			}

			now := time.Now().UnixNano()
			message := &model.OriginRow{
				Table:  name,
				Op:     op,
				Before: &before,
				After:  &after,
				Now:    now,
			}
			c.message <- message
			// log.Printf("send msg to c.message: %s", string(msg.Value))
		}

		sess.MarkMessage(msg, "")
	}
	if findErrorTable {
		log.Println("get error tables")
	}
	return nil
}
