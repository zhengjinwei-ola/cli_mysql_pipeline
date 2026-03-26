package room

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type mysqlAdapter struct {
	tableNames           []string
	pkField, sqlTemplate string
	orm                  orm.Ormer
}

type SQL string

func newMysqlAdapter(
	tableNames []string,
	pkField string,
	sqlTemplate SQL,
	orm orm.Ormer,
) *mysqlAdapter {
	return &mysqlAdapter{
		tableNames:  tableNames,
		pkField:     pkField,
		sqlTemplate: string(sqlTemplate),
		orm:         orm,
	}
}

func (a *mysqlAdapter) FindByIDs(ids []int64) ([]map[string]string, error) {
	strIDs := a.getStrIDs(ids)

	sql := fmt.Sprintf(a.sqlTemplate, strIDs)
	buf := make([]orm.Params, 0)
	_, err := a.orm.Raw(sql).Values(&buf)
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		for _, id := range ids {
			buf = append(buf, map[string]interface{}{
				"rid":  id,
				"sort": 0,
			})
		}
	}

	res := make([]map[string]string, 0, len(buf))
	for _, v := range buf {
		res = append(res, a.transformSingleRecord(v))
	}
	return res, nil
}

func (a *mysqlAdapter) GetNames() []string {
	return a.tableNames
}

func (a *mysqlAdapter) transformSingleRecord(record orm.Params) map[string]string {
	res := make(map[string]string)
	for k, v := range record {
		res[k] = fmt.Sprintf("%v", v)
	}
	return res
}

func (a *mysqlAdapter) getStrIDs(ids []int64) string {
	strIDs := make([]string, 0, len(ids))
	for _, v := range ids {
		strIDs = append(strIDs, strconv.FormatInt(v, 10))
	}
	return strings.Join(strIDs, ",")
}
