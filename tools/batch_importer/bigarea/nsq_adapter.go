package bigarea

import (
	"github.com/nsqio/go-nsq"
	"strconv"
)

type nsqAdapter struct {
	data       bigAreaMappings
	pkColName  string
	tableNames []string
}

func newNsqAdapter(pkColName string, tableNames []string) *nsqAdapter {
	return &nsqAdapter{
		pkColName:  pkColName,
		tableNames: tableNames,
	}
}

func (a *nsqAdapter) GetIDs(msg *nsq.Message) ([]int64, error) {
	if err := a.populateData(msg); err != nil {
		return nil, err
	}
	return a.data.getIDs(), nil
}

func (a *nsqAdapter) populateData(msg *nsq.Message) error {
	a.data = bigAreaMappingsFromNsq(msg)
	return nil
}

func (a *nsqAdapter) FindByIDs(ids []int64) ([]map[string]string, error) {
	bigAreas := a.data.getBigAreas()

	res := make([]map[string]string, 0, len(ids))
	for i := range bigAreas {
		res = append(res, make(map[string]string))
		res[i][a.pkColName] = strconv.FormatInt(ids[i], 10)
		res[i]["bigarea_id"] = strconv.FormatInt(bigAreas[i], 10)
	}
	return res, nil
}

func (a *nsqAdapter) GetNames() []string {
	return a.tableNames
}
