package common

import (
	"github.com/olachat/banban_server/common/serialize"
	"errors"
	"github.com/nsqio/go-nsq"
)

func NsqMsgToSlice(msg *nsq.Message) ([]interface{}, error) {
	buf, err := serialize.UnMarshal(msg.Body)
	if err != nil {
		return nil, err
	}

	bufArr, ok := buf.([]interface{})
	if !ok {
		return nil, errors.New("nsq msg is not an array")
	}
	if len(bufArr) == 0 {
		return nil, errors.New("empty ids received")
	}
	return bufArr, nil
}

func NsqMsgToMap(msg *nsq.Message) (map[string]interface{}, error) {
	buf, err := serialize.UnMarshal(msg.Body)
	if err != nil {
		return nil, err
	}

	bufArr, ok := buf.(map[string]interface{})
	if !ok {
		return nil, errors.New("nsq msg is not an array")
	}
	return bufArr, nil
}
