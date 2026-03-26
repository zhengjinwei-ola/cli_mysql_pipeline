package index

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Indexer struct {
	Name             string
	Mapping          map[string]interface{}
	NumberOfShards   int
	NumberOfReplicas int
}

func (this *Indexer) UpdateMapping() {
	post := map[string]interface{}{
		"properties": this.Mapping,
	}
	this.curl(
		"PUT",
		fmt.Sprintf("%s/_mapping/default", this.Name),
		post,
	)
}

func (this *Indexer) DropIndex() {
	this.curl("DELETE", this.Name, nil)
}

func (this *Indexer) CreateIndex() {
	numberOfShards := this.NumberOfShards
	numberOfReplicas := this.NumberOfReplicas
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
	this.curl("PUT", this.Name, post)
	this.curl(
		"PUT",
		fmt.Sprintf("%s/%s", this.Name, "_settings"),
		map[string]interface{}{
			"index.mapping.total_fields.limit": 10000,
			"index.mapper.dynamic":             false,
		},
	)
}

func (this *Indexer) curl(method string, path string, jsonData map[string]interface{}) error {
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
		log.Println("update mapping json \n", string(body))
	}

	output, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(output))
	if err != nil {
		panic(err)
	}
	return nil
}

func (this *Indexer) flowMain(message *model.Message) (string, error) {
	if message.Op == model.OpType_Write {
		return this.create(message.Id, message.Data)
	} else if message.Op == model.OpType_Update {
		return this.update(message.Id, message.Data)
	} else if message.Op == model.OpType_Delete {
		return this.remove(message.Id), nil
	}
	return "", fmt.Errorf("unknown ops type")
}

func (this *Indexer) flowJoin(message *model.Message) (string, error) {
	//删除数据，需要业务自己处理成update
	if message.Op == model.OpType_Write || message.Op == model.OpType_Update {
		return this.update(message.Id, message.Data)
	}
	return "", fmt.Errorf("unknown ops type")
}

func (this *Indexer) flowAppend(message *model.Message) (string, error) {
	//删除数据，需要业务自己处理成update
	if message.Op == model.OpType_Write || message.Op == model.OpType_Update {
		data := map[string]interface{}{}
		data[message.FlowField] = message.Data
		return this.update(message.Id, &data)
	}
	return "", fmt.Errorf("unknown ops type")
}

func (this *Indexer) create(docId int64, doc *map[string]interface{}) (string, error) {
	str, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(TplCreate, this.Name, docId, str), nil

}

func (this *Indexer) update(docId int64, value *map[string]interface{}) (string, error) {
	doc := *value
	_, ok := doc["script"]
	var data map[string]interface{}
	if ok {
		//是脚本
		data = doc
	} else {
		data = map[string]interface{}{
			"doc": doc,
		}
	}

	str, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(TplUpdate, this.Name, docId, string(str)), nil

}

func (this *Indexer) remove(docId int64) string {
	return fmt.Sprintf(TplRemove, this.Name, docId)
}
