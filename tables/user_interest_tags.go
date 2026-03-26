package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"
)

const MvaDelete string = `
def key = params.key;
ctx._source[key].removeIf(it -> it == params.value);
`
const MvaUpdate string = `
def key = params.key;
if(ctx._source[key].size() > 0){
	ctx._source[key].removeIf(it -> it == params.value);
}
ctx._source[key].add(params.value);
`

type OpRowUserInterestTags struct {
	OpRowTable
}

func (o OpRowUserInterestTags) Format(origin, before map[string]string, docField string,
	fields *[]model.EsField, op model.OpType) (int64, map[string]interface{}, model.OpType) {
	uid, err := strconv.ParseInt(origin[docField], 10, 64)
	cid, err1 := strconv.ParseInt(origin["cid"], 10, 64)
	if err != nil || err1 != nil || uid == 0 || cid == 0 {
		return 0, nil, op
	}
	switch op {
	case model.OpType_Write:
		return uid, o.wrap(cid, MvaUpdate), model.OpType_Update
	case model.OpType_Delete:
		return uid, o.wrap(cid, MvaDelete), model.OpType_Update
	case model.OpType_Update:
		return uid, o.wrap(cid, MvaDelete), model.OpType_Update
	}

	return 0, nil, op
}

func (o OpRowUserInterestTags) wrap(cid int64, script string) map[string]interface{} {
	return map[string]interface{}{
		"script": map[string]interface{}{
			"lang":   "painless",
			"inline": script,
			"params": map[string]interface{}{
				"key":   "interests",
				"value": cid,
			},
		},
	}
}
