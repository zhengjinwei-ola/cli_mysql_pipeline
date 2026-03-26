//go:generate mockgen -source=handler.go -destination=mocks/handler_mock.go -package=mocks
package common

import (
	"github.com/nsqio/go-nsq"
)

type UseCase interface {
	Import([]int64) error
}

type MsgAdapter interface {
	GetIDs(*nsq.Message) ([]int64, error)
}

type BatchImportHandler struct {
	uc UseCase
	ma MsgAdapter
}

func NewHandler(uc UseCase, ma MsgAdapter) *BatchImportHandler {
	return &BatchImportHandler{uc: uc, ma: ma}
}

func (h *BatchImportHandler) HandleMsg(msg *nsq.Message) error {
	ids, err := h.ma.GetIDs(msg)
	if err != nil {
		return err
	}

	return h.uc.Import(ids)
}
