package bigarea

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/common"
	"github.com/nsqio/go-nsq"
	"strconv"
)

type bigAreaMapping struct {
	id      int64
	bigArea int64
}

type bigAreaMappings []bigAreaMapping

func (m bigAreaMappings) getIDs() []int64 {
	res := make([]int64, 0, len(m))
	for _, item := range m {
		res = append(res, item.id)
	}
	return res
}

func (m bigAreaMappings) getBigAreas() []int64 {
	res := make([]int64, 0, len(m))
	for _, item := range m {
		res = append(res, item.bigArea)
	}
	return res
}

func bigAreaMappingsFromNsq(msg *nsq.Message) bigAreaMappings {
	mapping, _ := common.NsqMsgToMap(msg)
	res := make(bigAreaMappings, 0)
	for k, v := range mapping {
		id, _ := strconv.ParseInt(k, 10, 64)
		bigArea := v.(int)
		cur := bigAreaMapping{id: id, bigArea: int64(bigArea)}
		res = append(res, cur)
	}
	return res
}
