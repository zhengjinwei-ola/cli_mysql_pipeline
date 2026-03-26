package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/complement"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

var sqlGodSkill string = "select * from xs_user_god where uid = ?"

type OpRowUserProfileForGod struct {
	OpRowTable
}

func (o OpRowUserProfileForGod) Format(origin, before map[string]string, docField string,
	fields *[]model.EsField, op model.OpType) (int64, map[string]interface{}, model.OpType) {
	var err error
	var role int64 = 0
	var beforeRole int64 = 0
	changeToGod := false
	role, err = strconv.ParseInt(origin["role"], 10, 64)
	if err != nil || role < 2 {
		return 0, nil, op
	}

	if op == model.OpType_Update {
		beforeRole, err = strconv.ParseInt(before["role"], 10, 64)
		if err != nil {
			println(err.Error())
		}
		if role >= 2 && beforeRole < 2 {
			//用户成为大神
			op = model.OpType_Write
			changeToGod = true
		}
	}

	docId, data, op := o.OpRowTable.Format(origin, before, docField, fields, op)
	if docId > 0 && op == model.OpType_Write {
		//补齐默认数据
		data["room_rid"] = 0
		data["room_position"] = 0
		data["room_property"] = ""
		data["room_type"] = ""
		data["interests"] = []int64{}
		data["exposure"] = map[string]interface{}{
			"exposure_new":     100000000,
			"exposure_old":     100000000,
			"follow_new":       0,
			"income_yesterday": 0,
			"income_lastweek":  0,
			"income_total":     0,
			"is_peipei":        false,
			"follows":          []int64{},
			"interests":        []int64{},
			"friends_num":      0,
			"clicked_uids":     []int64{},
		}

		if interests, ok := origin["interests"]; ok {
			ids := []int64{}
			res := strings.Split(interests, ",")
			for _, cid := range res {
				v, err := strconv.ParseInt(cid, 10, 64)
				if err == nil {
					ids = append(ids, v)
				}
			}
			data["interests"] = ids
		}

		if changeToGod {
			//补齐技能数据
			res := []orm.Params{}
			_, err := model.MasterDb.Raw(
				sqlGodSkill,
				docId,
			).Values(&res)
			if err != nil {
				panic(err)
			}
			for i := 0; i < len(res); i++ {
				message := &model.OriginRow{
					Table:  "xianshi.xs_user_god",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  model.GetValueFromParams(res[i]),
					Full:   false,
				}
				complement.Instance.Add(message)
			}
		}
	}

	return docId, data, op
}
