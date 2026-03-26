//go:generate mockgen -source=use_case.go -destination=mocks/use_case_mock.go -package=mocks
package common

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"log"
)

type TableStorage interface {
	FindByIDs([]int64) ([]map[string]string, error)
	GetNames() []string
}

type BatchImportUseCase struct {
	index *index.Index
	table TableStorage
}

func NewBatchImportUseCase(
	index *index.Index, table TableStorage,
) *BatchImportUseCase {
	return &BatchImportUseCase{
		index: index,
		table: table,
	}
}

func (u *BatchImportUseCase) Import(ids []int64) error {
	records, err := u.table.FindByIDs(ids)
	if err != nil {
		return err
	}
	for i := range records {
		for _, tableName := range u.table.GetNames() {
			message := &model.OriginRow{
				Table:  tableName,
				Op:     model.OpType_Write,
				Before: &map[string]string{},
				After:  &records[i],
				Full:   true,
			}
			u.index.Receiver <- message
		}
	}
	log.Printf("successfully imported ids: %v for index: %v\n", ids, u.index.Name)
	return nil
}
