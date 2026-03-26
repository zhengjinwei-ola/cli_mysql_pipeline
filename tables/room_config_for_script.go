package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"
	"time"
)

const RoomConfigScript string = `
if(ctx._source.property == 'business' && (params.now - ctx._source.show_start) >= 3600){
	ctx._source.show_start = params.now;
}
`

type OpRowRoomConfigForScript struct {
	OpRowTable
}

func (this OpRowRoomConfigForScript) Format(
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
			data := this.getValue(uid)
			return rid, data, model.OpType_Update
		}
		break

	case model.OpType_Update:
		beforeUid, _ := strconv.ParseInt(before["uid"], 10, 64)
		beforePosition, _ := strconv.ParseInt(before["position"], 10, 64)
		if uid != beforeUid && (position == 0 || beforePosition == 0) {
			if uid > 0 {
				data := this.getValue(uid)
				return rid, data, model.OpType_Update
			}
		}

		break
	}

	return 0, nil, op
}

func (this OpRowRoomConfigForScript) getValue(
	uid int64,
) map[string]interface{} {
	return map[string]interface{}{
		"script": map[string]interface{}{
			"lang":   "painless",
			"inline": RoomConfigScript,
			"params": map[string]interface{}{
				"now": time.Now().Unix(),
			},
		},
	}
}
