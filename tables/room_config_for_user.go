package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"
)

var DefaultRoomConfig map[string]interface{} = map[string]interface{}{
	"room_rid":      0,
	"room_position": 0,
	"room_type":     "",
	"room_property": "",
}

const sqlRoomProperty string = "select rid, property, type from xs_chatroom where rid = ?"

type OpRowRoomConfigForUser struct {
	OpRowTable
}

func (o OpRowRoomConfigForUser) Format(origin, before map[string]string, docField string,
	fields *[]model.EsField, op model.OpType) (int64, map[string]interface{}, model.OpType) {
	docId, err := strconv.ParseInt(origin[docField], 10, 64)
	if err != nil {
		return 0, nil, op
	}
	switch op {
	case model.OpType_Write:
		if docId > 0 {
			data, err := o.format(origin, fields)
			if err != nil {
				return 0, nil, op
			}
			o.wrapData(&data, docId)
			return docId, data, model.OpType_Update
		}
	case model.OpType_Delete:
		if docId > 0 {
			return docId, DefaultRoomConfig, model.OpType_Update
		}
	case model.OpType_Update:
		beforeDocId, err := strconv.ParseInt(before[docField], 10, 64)
		if err != nil {
			return 0, nil, op
		}
		if beforeDocId != docId {
			//根据业务逻辑，rid和position不会变化
			//且beforeDocId和docId 不会同时大于0
			if docId > 0 {
				//上麦
				data, err := o.format(origin, fields)
				if err != nil {
					return 0, nil, op
				}
				o.wrapData(&data, docId)
				return docId, data, model.OpType_Update
			} else if beforeDocId > 0 {
				//下麦
				return beforeDocId, DefaultRoomConfig, model.OpType_Update
			}
		}
	}

	return 0, nil, op
}

func (o OpRowRoomConfigForUser) wrapData(data *map[string]interface{}, uid int64) {
	//根据房间rid, 查找房间属性
	value := *data
	rid, _ := value["room_rid"].(int64)
	room := model.XsChatroom{}
	err := model.Db.Raw(sqlRoomProperty, rid).QueryRow(&room)
	if err != nil {
		//log.Println("error get room", rid, err, value)
		value["room_type"] = ""
		value["room_property"] = ""
		value["room_rid"] = 0
		value["room_position"] = 0
	} else {
		value["room_type"] = room.Type
		value["room_property"] = room.Property
	}
}
