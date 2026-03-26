package tables

import (
	"fmt"
	"strconv"

	"github.com/kpango/glg"

	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
)

type OpRowRoom struct {
	OpRowTable
}

func (this OpRowRoom) Format(
	origin map[string]string,
	before map[string]string,
	docField string,
	fields *[]model.EsField,
	op model.OpType,
) (int64, map[string]interface{}, model.OpType) {
	docId, data, op := this.OpRowTable.Format(origin, before, docField, fields, op)
	fmt.Println("[debug-origin]", origin)
	fmt.Println("[debug-before]", before)
	fmt.Println("[debug-docId]", docId)
	fmt.Println("[debug-data]", data)
	fmt.Println("[debug-docField]", docField)
	fmt.Println("[debug-op]", op)

	if docId > 0 && op == model.OpType_Write {
		//补齐默认数据
		data["lock_end"] = 0
		data["emperor_time"] = 0
		data["reception_uid"] = 0
		data["fans_num"] = 0
		data["boy_num"] = 0
		data["girl_num"] = 0
		data["boss_uid"] = 0
		data["uname"] = ""
		data["utitle"] = 0
		data["sex"] = 0
		data["online_num"] = 0
		data["live_state"] = ""
		data["game_state_v2"] = 0
		data["game_state_create"] = 0
		data["game_state_update"] = 0
		data["link_mic_status"] = 0
		data["round_real"] = 0
		data["pk_state"] = 0
		data["red_packet_post_time"] = 0
		data["team_pk_state"] = 0
		data["multi_link_mic_status"] = 0
		data["multi_pk_status"] = 0

		data["base"] = map[string]interface{}{
			"new_male_stay_seconds":   0,
			"new_female_stay_seconds": 0,
			"old_male_stay_seconds":   0,
			"old_female_stay_seconds": 0,
			"new_male_pay_percent":    0.0,
			"new_female_pay_percent":  0.0,
			"old_male_pay_percent":    0.0,
			"old_female_pay_percent":  0.0,
		}

		// room condition defaults
		data["blocked"] = false
		// bigarea_id
		uid, _ := strconv.Atoi(fmt.Sprintf("%v", data["uid"]))
		if uid > 0 {
			bigareaId, err := model.GetBigAreaIdByUid(uid)
			if err != nil {
				glg.Error("OpRowRoom, fail get OpType_Write uid:", uid)
			} else {
				if bigareaId > 0 {
					data["bigarea_id"] = bigareaId
					glg.Debug("OpRowRoom union, bigareaId:", bigareaId, " data:", data)
				}
			}
		}
	}
	fmt.Println("[debug-data]", data)

	return docId, data, op
}
