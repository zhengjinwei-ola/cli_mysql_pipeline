package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"
	"strings"
)

const InterestTagsSql = "select cid from xs_user_interest_tags where uid = ?"

type InterestCid struct {
	Cid int64
}

type OpRowUserProfile struct {
	OpRowTable
}

func (this OpRowUserProfile) Format(
	origin map[string]string,
	before map[string]string,
	docField string,
	fields *[]model.EsField,
	op model.OpType,
) (int64, map[string]interface{}, model.OpType) {

	docId, data, op := this.OpRowTable.Format(origin, before, docField, fields, op)
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

		//game settings default
		data["sys_game_match_prompted_to"] = false
		data["mute_sys_game_prompts"] = false
		data["mute_stranger_game_invite"] = false
		data["mute_sys_game_prompts_till"] = 0
		data["mute_stranger_game_invite_till"] = 0
		data["viewing_page"] = -1
	}
	return docId, data, op
}
