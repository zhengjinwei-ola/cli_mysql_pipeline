package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"

	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
)

const (
	sqlRoomUsersAge string = `
			SELECT ccc.rid,
					MIN(YEAR(now()) - YEAR(uuu.birthday)) 'min_age', 
					MAX(YEAR(now()) - YEAR(uuu.birthday)) 'max_age', 
					AVG(YEAR(now()) - YEAR(uuu.birthday)) 'avg_age' 
			FROM xs_chatroom_config ccc 
			JOIN xs_user_profile uuu ON uuu.uid = ccc.uid 
			WHERE ccc.rid = ? AND ccc.uid > 0 AND uuu.birthday != 0 
			GROUP BY ccc.rid
`
)

type OpRowRoomConfigForRoomAge struct {
	OpRowTable
}

func (o OpRowRoomConfigForRoomAge) Format(
	origin map[string]string,
	before map[string]string,
	docField string,
	fields *[]model.EsField,
	op model.OpType,
) (int64, map[string]interface{}, model.OpType) {
	rid, err := strconv.ParseInt(origin[docField], 10, 64)
	if err != nil {
		return 0, nil, op
	}
	uid, _ := strconv.ParseInt(origin["uid"], 10, 64)

	switch op {
	case model.OpType_Write:
		if uid > 0 {
			data := o.getValue(rid)
			return rid, data, model.OpType_Update
		}
	case model.OpType_Delete:
	case model.OpType_Update:
		beforeUid, _ := strconv.ParseInt(before["uid"], 10, 64)
		if uid != beforeUid {
			data := o.getValue(rid)
			return rid, data, model.OpType_Update
		}
	}

	return 0, nil, op
}

func (o OpRowRoomConfigForRoomAge) getValue(rid int64) map[string]interface{} {

	buf := make([]orm.Params, 0)
	_, err := model.Db.Raw(sqlRoomUsersAge, rid).Values(&buf)
	if err != nil {
		println(err.Error())
	}

	data := map[string]interface{}{
		"min_age": 0,
		"max_age": 0,
		"avg_age": 0,
	}

	if len(buf) > 0 {
		o.convertToFloat(buf, data, "min_age")
		o.convertToFloat(buf, data, "max_age")
		o.convertToFloat(buf, data, "avg_age")
	}

	return data
}

func (o OpRowRoomConfigForRoomAge) convertToFloat(buf []orm.Params, data map[string]interface{}, col string) {
	switch v := buf[0][col].(type) {
	case int64, int32, int:
		data[col] = v.(float64)
	case float64, float32:
		data[col] = v
	case string:
		vv, _ := strconv.ParseFloat(v, 64)
		data[col] = vv
	default:
		log.Printf("unknown type %T for value %v", v, v)
	}
}
