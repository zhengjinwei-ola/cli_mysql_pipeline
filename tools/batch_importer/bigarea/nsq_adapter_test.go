package bigarea

import (
	"github.com/olachat/banban_server/common/serialize"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_phpAdapter_GetIDs(t *testing.T) {
	tests := []struct {
		name     string
		args     []map[int]int
		wantIDs  []int64
		wantRecs []map[string]string
		wantErr  bool
	}{
		{
			name: "happy case",
			args: []map[int]int{
				{891723: 1},
				{89787: 2},
			},
			wantIDs: []int64{891723, 89787},
			wantRecs: []map[string]string{
				{"uid": "891723", "bigarea_id": "1"},
				{"uid": "89787", "bigarea_id": "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &nsqAdapter{pkColName: "uid"}
			msg, err := serialize.Marshal(tt.args)
			nsqMsg := &nsq.Message{Body: msg}
			got, err := a.GetIDs(nsqMsg)
			assert.Equal(t, tt.wantIDs, got)
			assert.Equal(t, tt.wantErr, err != nil)

			got2, err := a.FindByIDs(tt.wantIDs)
			assert.Equal(t, tt.wantRecs, got2)
		})
	}
}
