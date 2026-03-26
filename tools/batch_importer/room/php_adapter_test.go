package room

import (
	"testing"

	"github.com/elliotchance/phpserialize"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

func Test_phpAdapter_GetIDs(t *testing.T) {
	tests := []struct {
		name    string
		param   []interface{}
		want    []int64
		wantErr bool
	}{
		{
			name:    "GIVEN mismatched type THEN error",
			param:   []interface{}{"1"},
			wantErr: true,
		},
		{
			name:    "GIVEN empty msg THEN error",
			param:   make([]interface{}, 0),
			wantErr: true,
		},
		{
			name:  "GIVEN non-empty msg THEN return list of int",
			param: []interface{}{8, 7, 4},
			want:  []int64{8, 7, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newPHPAdapter()
			msgBody, err := phpserialize.Marshal(&tt.param, nil)
			if err != nil {
				t.Error(err)
			}
			msg := nsq.NewMessage(nsq.MessageID{1}, msgBody)
			res, err := a.GetIDs(msg)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}
