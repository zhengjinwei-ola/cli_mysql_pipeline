package common

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/room/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestBatchImportUseCase_Import(t *testing.T) {
	type tableRes struct {
		data []map[string]string
		err  error
	}
	tests := []struct {
		name       string
		idx        *index.Index
		tableNames []string
		tableRes   tableRes
		ids        []int64
		wantErr    bool
		want       []*model.OriginRow
	}{
		{
			name:     "GIVEN error finding ids in table storage THEN error",
			idx:      &index.Index{Receiver: make(chan *model.OriginRow)},
			tableRes: tableRes{err: errors.New("some error")},
			ids:      make([]int64, 0),
			wantErr:  true,
			want:     make([]*model.OriginRow, 0),
		},
		{
			name:       "GIVEN no error finding ids in table storage THEN ok",
			idx:        &index.Index{Receiver: make(chan *model.OriginRow)},
			tableNames: []string{"test", "qer"},
			tableRes: tableRes{
				data: []map[string]string{
					{"id": "1"},
					{"id": "3"},
					{"id": "4"},
				},
			},
			ids: []int64{1, 3, 4},
			want: []*model.OriginRow{
				{
					Table:  "test",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  &map[string]string{"id": "1"},
					Full:   true,
				},
				{
					Table:  "test",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  &map[string]string{"id": "3"},
					Full:   true,
				},
				{
					Table:  "test",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  &map[string]string{"id": "4"},
					Full:   true,
				},
				{
					Table:  "qer",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  &map[string]string{"id": "1"},
					Full:   true,
				},
				{
					Table:  "qer",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  &map[string]string{"id": "3"},
					Full:   true,
				},
				{
					Table:  "qer",
					Op:     model.OpType_Write,
					Before: &map[string]string{},
					After:  &map[string]string{"id": "4"},
					Full:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockTable := mocks.NewMocktableStorage(ctrl)
			mockTable.EXPECT().FindByIDs(tt.ids).
				Return(tt.tableRes.data, tt.tableRes.err).Times(1)
			mockTable.EXPECT().GetNames().Return(tt.tableNames).AnyTimes()

			uc := NewBatchImportUseCase(tt.idx, mockTable)
			res := make([]*model.OriginRow, 0, 10)

			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				for msg := range tt.idx.Receiver {
					res = append(res, msg)
				}
				wg.Done()
			}()

			err := uc.Import(tt.ids)
			close(tt.idx.Receiver)
			wg.Wait()

			assert.Equal(t, tt.wantErr, err != nil)
			assert.ElementsMatch(t, tt.want, res)
		})
	}
}
