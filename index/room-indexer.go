package index

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RoomIndexer struct {
	Indexer
}

/**
 	UpdateRoomEs takes a rid and a SimpleOpRowBase to retrieve from mySQL to update ES

 	Note that it will follow op's FlowType to update existing index or create new index
	creating new index will fail when rid exist. update will fail when rid does not exist.
	TODO: upsert records
*/
func (this *RoomIndexer) UpdateRoomEs(rid int, op tables.SimpleOpRowBase) error {
	messages := op.RetrieveRows("rid", rid)
	output := make([]string, 0, len(messages))
	for _, row := range messages {

		if row.Full {
			op.Wrap.InitHookBefore(row.After)
		}
		docId, doc, opType := op.Wrap.Format(
			*row.After,
			*row.Before,
			op.DocField,
			op.Fields,
			row.Op,
		)
		if docId == 0 || doc == nil {
			continue
		}

		if row.Full {
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
			message.FlowField = op.FlowField(*row.After)
		}

		//根据flow和op写入数据
		var err error
		var o string
		switch message.Flow {
		case model.FlowType_Main:
			o, err = this.flowMain(message)

		case model.FlowType_Join:
			o, err = this.flowJoin(message)

		case model.FlowType_Append:
			o, err = this.flowAppend(message)

		default:
			err = fmt.Errorf("unknown message flow")
		}
		if err == nil {
			output = append(output, o)
		}

		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	fmt.Println(op.Name, "has", len(output), "records")
	fmt.Println(output)
	return bulkIngest(output)
}

func bulkIngest(msges []string) error {
	step := 500
	for i := 0; i < len(msges); i += step {
		end := i + step
		if end > len(msges) {
			end = len(msges)
		}
		x := strings.Join(msges[i:end], "\n")
		curl("PUT", "_bulk", x)
	}
	return nil
}

func curl(method string, path string, data string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", conf.EsHost, path), bytes.NewReader([]byte(data)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(conf.EsName) > 0 && len(conf.EsPass) > 0 {
		println("conf.EsName", conf.EsName)
		req.SetBasicAuth(conf.EsName, conf.EsPass)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res map[string]interface{}
	err = json.Unmarshal(output, &res)
	if err != nil {
		panic(err)
	}
	success := !res["errors"].(bool)
	if !success {
		fmt.Println("failed in some docs")
		println(string(output))
	}
	return nil
}
