package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/common/nsq"
	"github.com/olachat/banban_server/common/serialize"

	goNsq "github.com/nsqio/go-nsq"
)

var allowTable = map[string]bool{
	"xianshi.es_exposure":        true,
	"xianshi.es_update_lock_end": true,
	"xianshi.es_online_num":      true,
	"xianshi.xs_chatroom":        true,
	"xianshi.xs_circle_num":      true,
	"xianshi.xs_day_score":       true,
}

type handler func(*goNsq.Message) *model.OriginRow

//来自用户行为上报
func NewNsqMq(stop chan bool, message chan *model.OriginRow, h handler) {
	temp := strings.Split(conf.NsqAddr, ",")
	addrs := make([]string, 0)
	for _, address := range temp {
		addrs = append(addrs, strings.TrimSpace(address))
	}
	glg.Info("conf item nsq", conf.NsqAddr)

	//监听用户上报行为
	consumerMessage := nsq.NewNsqConsumer(
		nsq.Topic(conf.NsqTopic),
		nsq.Channel(conf.NsqChannel),
		func(msg *goNsq.Message) error {
			row := h(msg)
			if row != nil {
				message <- row
			}
			return nil
		},
	)
	err := consumerMessage.Connect(addrs)
	if err != nil {
		panic(err)
	}
	log.Println("nsq connect message ok")

	<-stop
	consumerMessage.Client.Stop()
}

func DefaultConverter(msg *goNsq.Message) *model.OriginRow {
	if len(msg.Body) == 0 {
		return nil
	}

	result, err := serialize.UnMarshal(msg.Body)
	if err != nil {
		glg.Error(err)
		return nil
	}

	//log.Printf("%T | %#v\n", result, result)
	res, ok := result.(map[string]interface{})

	// tools.IssueLog("JSON:",string(data))
	if !ok {
		glg.Errorf("error result data")
		return nil
	}

	data, ok := res["data"]
	if ok {
		if dataArr, ok := data.([]map[string]interface{}); ok {
			if len(dataArr) > 0 {
				dataMap := dataArr[0]
				if before, ok := dataMap["before"]; ok {
					res["before"] = before
				}
				if after, ok := dataMap["after"]; ok {
					res["after"] = after
				}
			}
		}
	}
	jsondata, err := json.Marshal(res)
	if err == nil {
		log.Printf("=== nsq handle, msg body : %s \n", string(jsondata))
	}

	row := model.OriginRow{}
	var opType string
	for key, value := range res {
		switch key {
		case "table":
			row.Table, ok = value.(string)
			if !ok {
				glg.Error("error result field table")
				return nil
			}
		case "type":
		case "op":
			opType, ok = value.(string)
			if !ok {
				glg.Error("error result field table")
				return nil
			}

			switch opType {
			case "write":
				row.Op = model.OpType_Write
			case "update":
				row.Op = model.OpType_Update
			case "delete":
				row.Op = model.OpType_Delete
			case "update_row":
				row.Op = model.OpType_Update_Row
			default:
				glg.Error("error result op", opType)
				return nil
			}
		case "before", "after", "update_row_data":
			val, ok := value.(map[string]interface{})
			if ok {
				value := model.GetValueFromParams(val)
				if key == "before" {
					row.Before = value
				} else if key == "update_row_data" {
					row.UpdateRowData = value
				} else {
					row.After = value
				}
			} else {
				glg.Error("error result after")
				return nil
			}
		case "doc_ids":
			val, ok := value.([]interface{})
			if !ok {
				glg.Errorf("error result field doc_names")
				continue
			}
			row.DocIds = &val
			break
		case "doc_field":
			row.DocField, ok = value.(string)
			if !ok {
				glg.Errorf("error result field doc_field")
				continue
			}
			break
		case "doc_names":
			val, ok := value.([]interface{})
			if !ok {
				glg.Errorf("error result field doc_names")
				continue
			}
			row.DocNames = &val
			break
		case "index_names":
			val, ok := value.([]interface{})
			if !ok {
				glg.Errorf("error result field index_names")
				continue
			}
			row.IndexNames = &val
			break
		}
	}

	if len(row.Table) == 0 || (row.After == nil && row.Op != model.OpType_Update_Row) {
		glg.Error("error result")
		return nil
	}

	if row.Op == model.OpType_Update && row.Before == nil {
		glg.Error("error result before")
		return nil
	}

	if row.Before == nil {
		row.Before = &map[string]string{}
	}
	log.Println("jinwei-debug", row.After)
	return &row
}

func TestConverter(msg *goNsq.Message) *model.OriginRow {
	if len(msg.Body) == 0 {
		return nil
	}

	result, err := serialize.UnMarshal(msg.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	res, ok := result.(map[string]interface{})
	if !ok {
		glg.Error("error result data")
		return nil
	}
	row := &model.OriginRow{}
	if row.Table, err = getTableName(res); err != nil {
		log.Println(err)
		return nil
	}
	if row.Op, err = getOp(res); err != nil {
		log.Println(err)
		return nil
	}
	if row.Before, row.After, err = populateData(res, row.Op); err != nil {
		log.Println(err)
		return nil
	}
	return row
}

func populateData(msg map[string]interface{}, op model.OpType) (*map[string]string, *map[string]string, error) {
	data, ok := msg["data"]
	if !ok {
		return nil, nil, errors.New("missing \"data\" in nsq message")
	}
	v, _ := data.([]interface{})
	vv, _ := v[0].(map[string]interface{})
	switch op {
	case model.OpType_Update:
		vvv, ok := vv["before"].(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("illy formatted msg: %T\n %v", data, data)
		}
		before := model.GetValueFromParams(vvv)
		vvv, ok = vv["after"].(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("illy formatted msg: %v", data)
		}
		after := model.GetValueFromParams(vvv)
		return before, after, nil
	case model.OpType_Write, model.OpType_Delete:
		after := model.GetValueFromParams(vv)
		return &map[string]string{}, after, nil
	default:
		return nil, nil, fmt.Errorf("unknown op: %v", op)
	}
}

func getOp(msg map[string]interface{}) (model.OpType, error) {
	t, ok := msg["type"]
	if !ok {
		return 0, errors.New("missing \"type\" in nsq msg")
	}
	switch t.(string) {
	case "write":
		return model.OpType_Write, nil
	case "update":
		return model.OpType_Update, nil
	case "delete":
		return model.OpType_Delete, nil
	default:
		return 0, fmt.Errorf("unknown type: %v", t.(string))
	}
}

func getTableName(msg map[string]interface{}) (string, error) {
	db, _ := msg["db"].(string)
	table, _ := msg["table"].(string)
	return db + "." + table, nil
}
