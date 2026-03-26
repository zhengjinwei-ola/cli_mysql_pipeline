package index

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"github.com/syyongx/php2go"
)

const (
	TplCreate = "{\"index\":{\"_index\":\"%s\",\"_type\":\"default\",\"_id\":\"%d\"}}\n%s\n"
	TplUpdate = "{\"update\":{\"_index\":\"%s\",\"_type\":\"default\",\"_id\":\"%d\"}}\n%s\n"
	TplRemove = "{\"delete\":{\"_index\":\"%s\",\"_type\":\"default\",\"_id\":\"%d\"}}\n"
)

type Index struct {
	Ops              []tables.OpRowBase //哪些表使用哪个类去操作
	Receiver         chan *model.OriginRow
	Name             string
	buffer           bytes.Buffer
	fp               *os.File
	num              int64
	count            int64
	execDocNum       int64
	stop             chan bool
	wg               *sync.WaitGroup
	opsMap           map[string][]tables.OpRowBase
	Mapping          map[string]interface{}
	NumberOfShards   int
	NumberOfReplicas int
	currentName      string
}

func (idx *Index) Init(wg *sync.WaitGroup) {
	idx.wg = wg
	idx.stop = make(chan bool)
	idx.opsMap = make(map[string][]tables.OpRowBase)
	for _, op := range idx.Ops {
		name := fmt.Sprintf("%s.%s", op.Db, op.Name)
		if _, ok := idx.opsMap[name]; !ok {
			idx.opsMap[name] = make([]tables.OpRowBase, 0)
		}
		idx.opsMap[name] = append(idx.opsMap[name], op)
	}
}

func (idx *Index) Full(autoDropIndex bool) {
	//尝试创建index
	if autoDropIndex {
		fmt.Println("drop start", idx.Name)
		os.Remove(idx.getLogName())
		idx.dropIndex()
	}

	idx.createIndex()
	idx.UpdateMapping()
	for _, op := range idx.Ops {
		idx.currentName = op.Name
		op.Init(idx.Receiver)
	}
}

func (idx *Index) UpdateMapping() {
	post := map[string]interface{}{
		"properties": idx.Mapping,
	}
	idx.curl(
		"PUT",
		fmt.Sprintf("%s/_mapping/default", idx.Name),
		post,
	)
}

func (idx *Index) dropIndex() {
	idx.curl("DELETE", idx.Name, nil)
}

func (idx *Index) createIndex() {
	numberOfShards := idx.NumberOfShards
	numberOfReplicas := idx.NumberOfReplicas
	if conf.IsDev {
		//ES 集群是三个节点
		numberOfShards = 1
		numberOfReplicas = 0
	}
	post := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   numberOfShards,
			"number_of_replicas": numberOfReplicas,
		},
		"mapping": map[string]interface{}{
			"default": map[string]interface{}{
				"properties": map[string]interface{}{},
			},
		},
	}
	idx.curl("PUT", idx.Name, post)
	idx.curl(
		"PUT",
		fmt.Sprintf("%s/%s", idx.Name, "_settings"),
		map[string]interface{}{
			"index.mapping.total_fields.limit": 10000,
			"index.mapper.dynamic":             false,
		},
	)
}

func (ix *Index) curl(method string, path string, jsonData map[string]interface{}) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", conf.EsHost, path), bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(conf.EsName) > 0 && len(conf.EsPass) > 0 {
		println("conf.EsName", conf.EsName)
		req.SetBasicAuth(conf.EsName, conf.EsPass)
	}

	resp, err := client.Do(req)

	if jsonData != nil {
		glg.Info("update mapping json \n", string(body))
	}

	output, _ := ioutil.ReadAll(resp.Body)
	glg.Info(string(output))
	if err != nil {
		panic(err)
	}
	return nil
}

func (idx *Index) Listen() {
	glg.Info("Listen Index:", idx.Name, ", listen")
	idx.flush(true)
	timer := time.NewTicker(time.Millisecond * 200)
	timerFlush := time.NewTicker(time.Second * 3)

	defer func() {
		timer.Stop()
		timerFlush.Stop()
	}()

	for {
		select {
		case msg := <-idx.Receiver:
			name := msg.Table
			glg.Infof("banban-es table:%s", name)
			glg.Infof("banban-es msgBefore:%+v", msg.Before)
			glg.Infof("banban-es msgAfter:%+v", msg.After)
			glg.Infof("banban-es msgOp:%+v", msg.Op)
			ops, ok := idx.opsMap[name]
			if !ok {
				panic(fmt.Errorf("Index %s Receiver Message %s", idx.Name, name))
			}

			for _, op := range ops {

				//tools.IssueLog("    DocName:[%s], op.Name:[%v]", idx.Name, op.Name)
				if msg.DocNames != nil && !php2go.InArray(idx.Name, *msg.DocNames) {
					//tools.IssueLog("        Skipping doc [%s] for table: [%s]", idx.Name, name)
					continue
				}

				if msg.IndexNames != nil && (op.IndexName == "" || !php2go.InArray(op.IndexName, *msg.IndexNames)) {
					//tools.IssueLog("        Skipping index [%s] for table: [%s],IndexNames:[%v]", op.IndexName, name, *msg.IndexNames)
					continue
				}
				// special for row update
				if msg.Op == model.OpType_Update_Row {

					// to populate the 'after' field based on op.SqlTemplate
					// RetrieveRows will resend originRow to this.receiver with msg.Op==opType_Write
					messages, err := op.RetrieveRow(msg)
					if err != nil {
						continue
					}

					for _, _message := range messages {
						// idx.Receiver <- _message
						// tools.IssueLog("        Looping index [%s][%s] for table: [%s], idx.Name:[%s], idx.currentName:[%s], name:[%s] ", counter, op.Name, op.IndexName, idx.Name, idx.currentName, name)

						idx.Execute(_message, op)
					}

					continue
				}
				log.Println("jjjj name msg", msg.After)
				idx.Execute(msg, op)

				//log.Println("used ", this.used(msg.Now))
			}
		case <-idx.stop:
			fmt.Println("wg done", idx.Name)
			idx.flush(true)
			idx.wg.Done()
			return
		case <-timer.C:
			//向磁盘刷新数据
			idx.flush(false)
		case <-timerFlush.C:
			idx.flush(true)
		}
	}
}

func (idx *Index) Execute(msg *model.OriginRow, op tables.OpRowBase) {
	if msg.Full {
		op.Wrap.InitHookBefore(msg.After)
	}
	docId, doc, opType := op.Wrap.Format(
		*msg.After,
		*msg.Before,
		op.DocField,
		op.Fields,
		msg.Op,
	)
	log.Println("jinweiDebug Table doc", msg.Table, doc)
	// do not continue
	if docId == 0 || doc == nil {
		return
	}

	if msg.Full {
		op.Wrap.InitHook(docId, &doc)
	}
	message := &model.Message{
		Id:   docId,
		Data: &doc,
		Op:   opType,
		Flow: op.Flow,
		Pk:   op.PkField,
	}
	if op.Flow == model.FlowType_Append {
		message.FlowField = op.FlowField(*msg.After)
	}

	bytes, _ := json.Marshal(message.Data)
	log.Println("doc", message.Id, len(string(bytes)))
	//根据flow和op写入数据
	switch message.Flow {
	case model.FlowType_Main:
		log.Println("jinwei-debug flow type main")
		idx.flowMain(message, op.UseUpsert)
	case model.FlowType_Join:
		idx.flowJoin(message, op.UseUpsert)
	case model.FlowType_Append:
		idx.flowAppend(message)
	}
}

func (idx *Index) Stop() {
	idx.stop <- true
}

func (idx *Index) used(before int64) float64 {
	return float64(time.Now().UnixNano()-before) / 1000 / 1000
}

func (idx *Index) flush(force bool) {
	idx.count++
	if idx.buffer.Len() > 0 {
		idx.flushToDisk()
	}
	//刷新到ES
	if force || idx.num >= 10000 {
		idx.num = 0
		idx.count = 0
		for i := 0; i < 30; i++ {
			err := idx.exec()
			if err == nil {
				return
			} else {
				glg.Errorf("exec error with ", err, "try again")
				time.Sleep(time.Second * 5)
			}
		}
	}
}

func (idx *Index) exec() error {
	if idx.fp == nil {
		return nil
	}

	idx.num = 0
	info, err := idx.fp.Stat()
	if err != nil {
		glg.Error(err)
		return nil
	}
	if info.Size() == 0 {
		return nil
	}
	now := time.Now().UnixNano()

	glg.Infof("[index.go][311]  exec ES BEGIN, Name:[%s], currentName:[%s], execDocNum:[%d], log:[%s]", idx.Name, idx.currentName, idx.execDocNum, idx.getLogName())

	args := []string{}
	if len(conf.EsName) > 0 && len(conf.EsPass) > 0 {
		args = append(args, "-u", fmt.Sprintf("%s:%s", conf.EsName, conf.EsPass))
	}
	args = append(args, "-o", "/dev/null")
	args = append(args, "-XPUT", fmt.Sprintf("%s/_bulk", conf.EsHost))
	args = append(args, "-H", "Content-Type: application/json")
	args = append(args, "--data-binary", fmt.Sprintf("@%s", idx.getLogName()))
	cmd := exec.Command(
		"curl",
		args...,
	)
	glg.Infof("%+v", args)
	_, err = cmd.CombinedOutput()
	if err != nil {
		time.Sleep(time.Second * 3)
		return err
	}

	idx.fp.Truncate(0)

	glg.Infof("[index.go][332]  exec ES END, Name:[%s], currentName:[%s], execDocNum:[%d], time_used:[%f ms]", idx.Name, idx.currentName, idx.execDocNum, float64(time.Now().UnixNano()-now)/1000/1000)

	idx.execDocNum = 0
	return nil
}

func (idx *Index) flowMain(message *model.Message, upsert bool) {
	log.Println("jinwei-debug flow type main op:", message.Op)
	switch message.Op {
	case model.OpType_Write:
		idx.create(message.Id, message.Data, upsert)
	case model.OpType_Update:
		idx.update(message.Id, message.Data, upsert)
	case model.OpType_Delete:
		idx.remove(message.Id)
	}
}

func (idx *Index) flowJoin(message *model.Message, upsert bool) {
	//删除数据，需要业务自己处理成update
	if message.Op == model.OpType_Write || message.Op == model.OpType_Update {
		idx.update(message.Id, message.Data, upsert)
	}
}

func (idx *Index) flowAppend(message *model.Message) {
	//删除数据，需要业务自己处理成update
	if message.Op == model.OpType_Write || message.Op == model.OpType_Update {
		data := map[string]interface{}{}
		data[message.FlowField] = message.Data
		idx.update(message.Id, &data, false)
	}
}

func (idx *Index) create(docId int64, doc *map[string]interface{}, upsert bool) {
	str, err := json.Marshal(doc)
	if err == nil {
		if upsert {
			data := map[string]interface{}{
				"doc":           doc,
				"doc_as_upsert": true,
			}
			str, err := json.Marshal(data)
			if err == nil {
				idx.execDocNum++
				idx.write(fmt.Sprintf(TplUpdate, idx.Name, docId, string(str)))
			}
		} else {
			idx.write(fmt.Sprintf(TplCreate, idx.Name, docId, str))
		}
	}
}

func (idx *Index) update(docId int64, value *map[string]interface{}, upsert bool) {
	doc := *value
	_, ok := doc["script"]
	log.Println("jinwei-debug flow type main value:", doc)
	var data map[string]interface{}
	if ok {
		//是脚本
		data = doc
	} else {
		data = map[string]interface{}{
			"doc": doc,
		}
		if upsert {
			data["doc_as_upsert"] = true
		}
	}
	log.Println("jinwei-debug flow type main data:", data)
	str, err := json.Marshal(data)
	if err == nil {
		idx.execDocNum++
		idx.write(fmt.Sprintf(TplUpdate, idx.Name, docId, string(str)))
	}
}

func (idx *Index) remove(docId int64) {
	idx.write(fmt.Sprintf(TplRemove, idx.Name, docId))
}

func (idx *Index) write(line string) {
	idx.num++
	len, err := idx.buffer.WriteString(line)
	if err != nil {
		glg.Info("write", len, err)
	}
}

func (idx *Index) getLogName() string {
	return fmt.Sprintf("/tmp/es_%s.log", idx.Name)
}

func (idx *Index) flushToDisk() {
	if idx.fp == nil {
		var err error
		idx.fp, err = os.OpenFile(
			idx.getLogName(),
			os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666,
		)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	_, err := idx.fp.Write(idx.buffer.Bytes())
	if err != nil {
		glg.Error(err)
	}
	idx.buffer.Reset()
}
