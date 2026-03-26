package common

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/room/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

func TestBatchImportHandler_HandleMsg(t *testing.T) {
	type maRes struct {
		ids []int64
		err error
	}
	tests := []struct {
		name    string
		maRes   maRes
		ucErr   error
		arg     *nsq.Message
		wantErr bool
	}{
		{
			name:    "GIVEN error in msg adapter THEN error",
			maRes:   maRes{err: errors.New("some error")},
			arg:     &nsq.Message{},
			wantErr: true,
		},
		{
			name:    "GIVEN error in uc THEN error",
			maRes:   maRes{ids: []int64{1, 3, 5}},
			ucErr:   errors.New("uc error"),
			arg:     &nsq.Message{},
			wantErr: true,
		},
		{
			name:  "GIVEN no error in msg adapter and uc THEN ok",
			maRes: maRes{ids: []int64{87, 8, 1}},
			arg:   &nsq.Message{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			ma := mocks.NewMockmsgAdapter(ctrl)
			ma.EXPECT().GetIDs(tt.arg).Return(tt.maRes.ids, tt.maRes.err).Times(1)

			uc := mocks.NewMockuseCase(ctrl)
			uc.EXPECT().Import(tt.maRes.ids).Return(tt.ucErr).AnyTimes()

			h := NewHandler(uc, ma)
			err := h.HandleMsg(tt.arg)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
