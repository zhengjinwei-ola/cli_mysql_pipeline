package room

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/common"
	"fmt"

	"github.com/nsqio/go-nsq"
)

type phpAdapter struct {
}

func newPHPAdapter() *phpAdapter {
	return new(phpAdapter)
}

func (a *phpAdapter) GetIDs(msg *nsq.Message) ([]int64, error) {
	bufArr, err := common.NsqMsgToSlice(msg)
	if err != nil {
		return nil, err
	}

	res := make([]int64, 0, len(bufArr))
	for _, v := range bufArr {
		switch v := v.(type) {
		case int:
			res = append(res, int64(v))
		case int64:
			res = append(res, v)
		case float64:
			res = append(res, int64(v))
		default:
			return nil, fmt.Errorf("unable to convert %v (type: %T) to int64 type", v, v)
		}
	}
	return res, nil
}
