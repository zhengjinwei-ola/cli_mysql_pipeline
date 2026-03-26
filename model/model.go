package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var Db orm.Ormer
var MasterDb orm.Ormer
var BbuMasterDb orm.Ormer

func InitDb() {
	orm.RegisterDataBase("default", "mysql", conf.MysqlAddr, 1, 5)
	orm.RegisterDataBase("master", "mysql", conf.MysqlMasterAddr, 1, 5)

	orm.RegisterModel(
		new(XsChatroom),
	)
	orm.Debug = conf.MysqlDebug

	Db = orm.NewOrm()
	MasterDb = orm.NewOrm()
	err := MasterDb.Using("master")
	if err != nil {
		panic(err)
	}
}

// 数据表主键最大，最小
type TableProfile struct {
	Max int64
	Min int64
}

// 获取房间属性
type XsChatroom struct {
	Rid      int32 `orm:"pk"`
	Property string
	Type     string
}

type OverseaXsChatroom struct {
	Rid      int32 `orm:"pk"`
	Name     string
	Game     string
	Property string
	Types    string
	Type     string
	Weight   int32
}

type XsRoomConfigProfile struct {
	Uid   int32 `orm:"pk"`
	Name  string
	Icon  string
	Title int32
	Sex   int
}

type XsUserFriendCount struct {
	Total int32
}

type XsCategory struct {
	Cid int32 `orm:"pk"`
}

type Message struct {
	Id        int64
	Op        OpType
	Flow      FlowType
	FlowField string
	Pk        string
	Data      *map[string]interface{}
}

type OriginRow struct {
	Table         string
	DocField      string
	DocNames      *[]interface{}
	IndexNames    *[]interface{}
	Op            OpType
	Before        *map[string]string
	After         *map[string]string
	UpdateRowData *map[string]string
	DocIds        *[]interface{}
	Now           int64
	Full          bool
}

type EsField struct {
	Name   string
	Target string
	Type   EsType
	Func   func(origin map[string]string) interface{}
}

type EsGeo struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func (row *OriginRow) UnmarshalJSON(bodyRaw []byte) error {

	body := make(map[string]interface{})
	_ = json.Unmarshal(bodyRaw, &body)

	var opType string
	var ok bool
	for key, value := range body {
		switch key {
		case "table":
			row.Table, ok = value.(string)
			if !ok {
				glg.Errorf("error result field table")
				return nil
			}
			break
		case "op":
			opType, ok = value.(string)
			if !ok {
				glg.Errorf("error result field table")
				return nil
			}
			if opType == "write" {
				row.Op = OpType_Write
			} else if opType == "update" {
				row.Op = OpType_Update
			} else if opType == "delete" {
				row.Op = OpType_Delete
			} else if opType == "update_row" {
				row.Op = OpType_Update_Row
				// log.Printf("NSQ msg received: update row table fields %s: %s\n", res["table"].(string), res["before"].(map[string]interface{})["docId"])
			} else {
				glg.Errorf("[model.go]error result field op:%s", opType)
				return nil
			}
			break
		case "before", "after", "update_row_data":
			val, ok := value.(map[string]interface{})
			if ok {
				value := GetValueFromParams(val)
				if key == "before" {
					row.Before = value
				}
				if key == "update_row_data" {
					row.UpdateRowData = value
				} else {
					row.After = value
				}
			} else {
				return errors.New("error result field =>" + key)
			}
			break
		case "doc_ids":
			val, ok := value.([]interface{})
			if !ok {
				glg.Errorf("[model.go]error result field doc_names")
				continue
			}
			row.DocIds = &val
			break
		case "doc_field":
			row.DocField, ok = value.(string)
			if !ok {
				glg.Errorf("[model.go]error result field doc_field")
				continue
			}
			break
		case "doc_names":
			val, ok := value.([]interface{})
			if !ok {
				glg.Errorf("[model.go]error result field doc_names")
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
	return nil
}

func (row *OriginRow) GetOneIntoDocIds() (*[]int64, error) {
	if row.UpdateRowData == nil {
		return nil, errors.New("update_row_data doesn't exists")
	}

	docId := (*row.UpdateRowData)["doc_id"]

	if i, err := strconv.Atoi(docId); err == nil {
		var docIds []int64
		docIds = append(docIds, int64(i))
		return &docIds, nil
	} else {
		return nil, errors.New("failed to cast into int64")
	}
}

func (row *OriginRow) GetDocIds() (*[]int64, error) {
	if row.DocIds == nil {
		return nil, errors.New("no doc_ids defined")
	}

	var docIds []int64

	for _, did := range *row.DocIds {
		if id, err := getId(&did); err == nil {
			docIds = append(docIds, id)
		}
	}

	return &docIds, nil
}

func (row *OriginRow) GetString(key string) string {
	return (*row.UpdateRowData)[key]
}

func getId(v *interface{}) (int64, error) {

	switch (*v).(type) {
	case int:
		return int64((*v).(int)), nil
	case int64:
		return (*v).(int64), nil
	case float64:
		return int64((*v).(float64)), nil
	case string:
		if i, err := strconv.Atoi((*v).(string)); err == nil {
			return int64(i), nil
		} else {
			return 0, err
		}

	default:
		return 0, fmt.Errorf("unable to convert %v (type: %T) to int64 type", v, v)
	}
}
