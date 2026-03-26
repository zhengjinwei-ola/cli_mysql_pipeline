package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"

	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
)

var DefaultRoomConfigForRoom map[string]interface{} = map[string]interface{}{
	"reception_uid": 0,
	"icon":          "",
	"fans_num":      0,
	"uname":         "",
	"utitle":        0,
	"sex":           0,
}

const (
	sqlProfile      string = "select uid, name, title, icon, sex from xs_user_profile where uid = ?"
	sqlFriendCount  string = "select count(*) as total from xs_user_friend where `to` = ?"
 )

type OpRowRoomConfigForRoom struct {
	OpRowTable
}

func (o OpRowRoomConfigForRoom) Format(
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
	position, _ := strconv.ParseInt(origin["position"], 10, 64)

	switch op {
	case model.OpType_Write:
		if uid > 0 && position == 0 {
			data := o.getValue(uid, rid)
			return rid, data, model.OpType_Update
		}
	case model.OpType_Delete:
		if uid > 0 && position == 0 {
			return rid, DefaultRoomConfigForRoom, model.OpType_Update
		}
	case model.OpType_Update:
		beforeUid, _ := strconv.ParseInt(before["uid"], 10, 64)
		beforePosition, _ := strconv.ParseInt(before["position"], 10, 64)
		if uid != beforeUid && (position == 0 || beforePosition == 0) {
			if uid > 0 {
				data := o.getValue(uid, rid)
				return rid, data, model.OpType_Update
			} else {
				return rid, DefaultRoomConfigForRoom, model.OpType_Update
			}
		}
	}

	return 0, nil, op
}

func (o OpRowRoomConfigForRoom) getValue(uid, rid int64) map[string]interface{} {
	profile := model.XsRoomConfigProfile{}
	err := model.Db.Raw(sqlProfile, uid).QueryRow(&profile)
	data := map[string]interface{}{
		"reception_uid": uid,
		"sex":           0,
		"fans_num":      0,
	}
	if err != nil {
		data["utitle"] = 0
		data["sex"] = 0
		data["uname"] = ""
		data["icon"] = ""
	} else {
		data["utitle"] = profile.Title
		data["sex"] = profile.Sex
		data["uname"] = profile.Name
		data["icon"] = profile.Icon
	}

	count := model.XsUserFriendCount{}
	err = model.Db.Raw(sqlFriendCount, uid).QueryRow(&count)
	if err != nil {
		data["fans_num"] = count.Total
	}
	return data
}

func (o OpRowRoomConfigForRoom) convertToInt(buf []orm.Params, data map[string]interface{}, col string) {
	switch v := buf[0][col].(type) {
	case int64, int:
		data[col] = v
	case string:
		vv, _ := strconv.ParseInt(v, 10, 64)
		data[col] = vv
	default:
		log.Printf("unknown type %T for value %v", v, v)
	}
}
