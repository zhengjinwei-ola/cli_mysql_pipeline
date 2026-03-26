package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/common/serialize"
	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type OpRowRoomBigArea struct {
	OpRowTable
	nsqProducer *nsq.Producer
	topic       string
}

func NewOpRowRoomBigArea(topic string) *OpRowRoomBigArea {
	cfg := nsq.NewConfig()
	producer, _ := nsq.NewProducer(conf.NsqAddr, cfg)
	return &OpRowRoomBigArea{topic: topic, nsqProducer: producer}
}

func (o *OpRowRoomBigArea) Format(
	origin map[string]string,
	before map[string]string,
	docField string,
	fields *[]model.EsField,
	op model.OpType,
) (int64, map[string]interface{}, model.OpType) {
	if _, ok := origin[docField]; ok {
		return o.OpRowTable.Format(origin, before, docField, fields, op)
	}
	o.populateRIDs(origin)
	return 0, nil, 0
}

func (o OpRowRoomBigArea) populateRIDs(origin map[string]string) {
	uidRaw, ok := origin["uid"]
	if !ok {
		log.Errorf("doc doesn't have uid: %v", origin)
		return
	}
	uid, _ := strconv.ParseInt(uidRaw, 10, 64)

	bigAreaID, ok := origin["bigarea_id"]
	if !ok {
		log.Errorf("doc doesn't have bigarea: %v", origin)
		return
	}
	for _, rid := range o.getRIDsByUID(uid) {
		data := map[string]string{
			"rid":        rid,
			"bigarea_id": bigAreaID,
		}
		body, _ := serialize.Marshal(data)
		o.nsqProducer.Publish(o.topic, body)
	}
}

func (o OpRowRoomBigArea) getRIDsByUID(uid int64) []string {
	var rids []string
	_, _ = model.Db.Raw("select rid from xianshi.xs_chatroom where uid = ?", uid).QueryRows(&rids)
	return rids
}
