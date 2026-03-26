package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/kpango/glg"
	"strconv"
)

const exposureNewUpdate string = `
if(ctx._source.exposure.exposure_new == 0){
	ctx._source.exposure.exposure_new += (params.num + 100000000 )
}else{
	ctx._source.exposure.exposure_new += params.num;
}
`

const exposureOldUpdate string = `
if(ctx._source.exposure.exposure_old == 0){
	ctx._source.exposure.exposure_old += (params.num + 100000000 )
}else{
	ctx._source.exposure.exposure_old += params.num; 
}
`

type OpRowEsExposure struct {
	OpRowTable
}

func (o OpRowEsExposure) Format(origin, before map[string]string, docField string,
	fields *[]model.EsField, op model.OpType) (int64, map[string]interface{}, model.OpType) {
	uid, err := strconv.ParseInt(origin[docField], 10, 64)
	num, err1 := strconv.ParseInt(origin["num"], 10, 64)
	if uid == 0 || err != nil || err1 != nil {
		glg.Error("OpRowEsExposure Format break", origin, err, err1)
		return 0, nil, op
	}
	var script string
	if origin["type"] == "exposure_new" {
		script = exposureNewUpdate
	} else {
		script = exposureOldUpdate
	}
	return uid, o.wrap(uid, num, script), model.OpType_Update
}

func (o OpRowEsExposure) wrap(uid int64, num int64, script string) map[string]interface{} {
	return map[string]interface{}{
		"script": map[string]interface{}{
			"lang":   "painless",
			"inline": script,
			"params": map[string]interface{}{
				"id":  uid,
				"num": num,
			},
		},
	}
}
