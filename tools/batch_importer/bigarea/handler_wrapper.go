//go:generate mockgen -source=handler_wrapper.go -destination=mocks/handler_wrapper_mock.go -package=mocks
package bigarea

import (
	"github.com/olachat/banban_server/common/nsq"
	"github.com/olachat/banban_server/common/serialize"
	"encoding/json"
	"fmt"
	goNsq "github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

type nsqMsg struct {
	Type string      `json:"type"`
	Data map[int]int `json:"data"`
}

type iBatchImportHandler interface {
	HandleMsg(msg *goNsq.Message) error
}

func WrapBigAreaHandlers(hRoom, hUser iBatchImportHandler) nsq.NsqHandleMessage {
	return func(msg *goNsq.Message) error {
		body, err := msgBody(msg)
		if err != nil {
			log.Error(err)
			return nil
		}
		newMsg := body.getData()
		switch body.Type {
		case "user.bigarea.batch":
			return hUser.HandleMsg(newMsg)
		case "room.bigarea.batch":
			return hRoom.HandleMsg(newMsg)
		default:
			log.Errorf("unknown command: {%v}", body.Type)
			return nil
		}
	}
}

func (m *nsqMsg) getData() *goNsq.Message {
	buf, _ := serialize.Marshal(m.Data)
	newMsg := &goNsq.Message{Body: buf}
	return newMsg
}

func msgBody(msg *goNsq.Message) (*nsqMsg, error) {
	buf, err := serialize.UnMarshal(msg.Body)
	if err != nil {
		return nil, err
	}
	bufv, _ := buf.(map[string]interface{})
	t, _ := bufv["type"].(string)
	switch t {
	case "user.bigarea.batch":
	case "room.bigarea.batch":
	default:
		return nil, fmt.Errorf("unknown command %v", t)
	}
	jsonBuf, err := json.Marshal(buf)
	if err != nil {
		return nil, err
	}
	res := new(nsqMsg)
	if err = json.Unmarshal(jsonBuf, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
