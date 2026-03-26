package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"strconv"
)

const MvaUpd string = `
if (ctx._source.gamecard_cids == null) { ctx._source.gamecard_cids = new ArrayList(); }
if (!ctx._source.gamecard_cids.contains(params.cid)) { ctx._source.gamecard_cids.add(params.cid); }
`
const MvaDel string = `if (ctx._source.gamecard_cids != null && ctx._source.gamecard_cids.contains(params.cid)) { ctx._source.gamecard_cids.remove(ctx._source.gamecard_cids.indexOf(params.cid)) }`

type OpRowUserGameCard struct {
	OpRowTable
}

func (o OpRowUserGameCard) Format(origin, before map[string]string, docField string,
	fields *[]model.EsField, op model.OpType) (int64, map[string]interface{}, model.OpType) {
	uid, err := strconv.ParseInt(origin[docField], 10, 64)
	cid, err1 := strconv.ParseInt(origin["cid"], 10, 64)
	if err != nil || err1 != nil || uid == 0 || cid == 0 {
		return 0, nil, op
	}
	switch op {
	case model.OpType_Write:
		return uid, o.wrap(cid, MvaUpd), model.OpType_Update
	case model.OpType_Delete:
		return uid, o.wrap(cid, MvaDel), model.OpType_Update
	case model.OpType_Update:
		return uid, o.wrap(cid, MvaDel), model.OpType_Update
	}

	return 0, nil, op
}

func (o OpRowUserGameCard) wrap(cid int64, script string) map[string]interface{} {
	return map[string]interface{}{
		"script": map[string]interface{}{
			"lang":   "painless",
			"inline": script,
			"params": map[string]interface{}{
				"cid": cid,
			},
		},
	}
}
