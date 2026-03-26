package tables

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools"
	"github.com/syyongx/php2go"

	"github.com/astaxie/beego/orm"
)

type OpRowWrap interface {
	Format(
		origin map[string]string,
		before map[string]string,
		docField string,
		fields *[]model.EsField,
		op model.OpType,
	) (int64, map[string]interface{}, model.OpType)

	InitHook(docId int64, origin *map[string]interface{})
	InitHookBefore(origin *map[string]string)
}

type OpRowBase struct {
	IndexName string
	Name      string
	Db        string
	/*
		文档依赖ID
	*/
	DocField string
	/*
		主键ID
	*/
	PkField            string
	SqlTemplate        string
	BatchSqlTemplate   string
	MinMaxSqlTemplate  string
	GappingSqlTemplate string
	/*
		如何格式数据
	*/
	Fields *[]model.EsField
	/*
		操作类型
	*/
	Flow      model.FlowType
	FlowField func(map[string]string) string
	Wrap      OpRowWrap
	/*
		Special handling for index creation even if it doesnt exists. Be careful when using it as it will corrupt the date
		if upsert happens not in sync.
	*/
	UseUpsert bool
}

//全量初始化时
func (o OpRowBase) Init(receiver chan *model.OriginRow) {
	if len(o.PkField) == 0 {
		return
	}
	//获取主键的最大值和最小值
	profile := model.TableProfile{}

	sql := fmt.Sprintf(
		"select max(%s) as max, min(%s) as min from %s",
		o.PkField,
		o.PkField,
		o.Name,
	)

	if len(o.MinMaxSqlTemplate) > 0 {
		sql = o.MinMaxSqlTemplate
	}

	err := model.Db.Raw(sql).QueryRow(&profile)
	if err != nil {
		glg.Error(sql)
		glg.Error(err)
		//panic(err)
		err = nil

		// due to migration and isolate from PT, hence some table row record might not be around while still receiving
		// NSQ packets, direct return at the moment
		return
	}

	if profile.Max <= 0 {
		return
	}
	var min = profile.Min
	var max = profile.Max + 1
	var step int64 = 1000
	var fullName = fmt.Sprintf("%s.%s", o.Db, o.Name)

	//根据主键范围，循环查询数据
	sql = fmt.Sprintf(
		"select * from %s where %s >= ? and %s < ?",
		o.Name,
		o.PkField,
		o.PkField,
	)
	if len(o.SqlTemplate) > 0 {
		sql = o.SqlTemplate
	}

	numOfQuestions := strings.Count(sql, "?")

	for begin := min; begin <= max; begin += step {
		res := []orm.Params{}
		var argssss []interface{}
		for i := 0; i < numOfQuestions; i++ {
			if i%2 == 0 {
				argssss = append(argssss, begin)
			} else {
				argssss = append(argssss, begin+step)
			}
		}
		_, err := model.Db.Raw(
			sql,
			argssss...,
		).Values(&res)

		if err != nil {
			glg.Error(sql)
			glg.Error(err)
			//panic(err)
			err = nil

			// due to migration and isolate from PT, hence some table row record might not be around while still receiving
			// NSQ packets, direct return at the moment
			return
		}
		for i := 0; i < len(res); i++ {
			message := &model.OriginRow{
				Table:  fullName,
				Op:     model.OpType_Write,
				Before: &map[string]string{},
				After:  model.GetValueFromParams(res[i]),
				Full:   true,
			}
			receiver <- message
		}

		if len(res) == 0 {
			//查询下一个节点
			//解决海外用户UID跳跃过大的问题
			nextSql := fmt.Sprintf(
				"select %s as min from %s where %s >= ? order by %s asc limit 1",
				o.PkField,
				o.Name,
				o.PkField,
				o.PkField,
			)

			// has specific statement
			if len(o.GappingSqlTemplate) > 0 {
				nextSql = o.GappingSqlTemplate
			}

			numOfMarks := strings.Count(nextSql, "?")
			var nextArgs []interface{}
			for i := 0; i < numOfMarks; i++ {
				nextArgs = append(nextArgs, begin+step)
			}

			next := model.TableProfile{}
			err := model.Db.Raw(nextSql, nextArgs...).QueryRow(&next)
			if err != nil {
				if err == orm.ErrNoRows {
					return
				}

				glg.Error(nextSql)
				glg.Error(err)
				//panic(err)
				err = nil

				// due to migration and isolate from PT, hence some table row record might not be around while still receiving
				// NSQ packets, direct return at the moment
				return
			} else {
				begin = next.Min - step
			}
		}
	}
}

func (op OpRowBase) RetrieveRow(originRow *model.OriginRow) ([]*model.OriginRow, error) {

	var fullName = fmt.Sprintf("%s.%s", op.Db, op.Name)
	docIds, err := originRow.GetOneIntoDocIds()

	//
	if err != nil {
		// try if doc_ids defined
		docIds, err = originRow.GetDocIds()
		if err != nil {
			return nil, errors.New("No ")
		}
	}

	sql := ""
	if len(op.BatchSqlTemplate) > 0 {
		sql = op.BatchSqlTemplate
	} else {
		sql = fmt.Sprintf(
			"select * from %s where %s IN (?) ",
			op.Name,
			op.DocField,
		)
	}

	marks := make([]string, 0, len(*docIds))
	tmpDocIds := make([]interface{}, 0, len(*docIds))

	for _, tmpDocId := range *docIds {
		marks = append(marks, "?")
		tmpDocIds = append(tmpDocIds, tmpDocId)
	}

	sql = strings.Replace(sql, "?", php2go.Implode(",", marks), -1)

	//tools.IssueLog("RetrieveRow for %s with SQL:%s", fullName, sql)
	//glg.Info("[op_row.go][208]", sql)

	var res []orm.Params
	_, err = model.Db.Raw(
		sql,
		tmpDocIds...,
	).Values(&res)

	if err != nil {
		panic("fail to retrieve data for " + fullName + ". err: " + err.Error())
	}

	if len(res) == 0 {
		//glg.Warnf("[op_row.go]no records found for SQL:%s ||| %v", sql, docIds)
		return nil, errors.New("no records found for SQL")
	}
	//tools.IssueLog("    [%s][%d] records found for SQL:%s",fullName, len(res), sql)
	messages := make([]*model.OriginRow, 0, len(res))
	for _, raw := range res {
		after := model.GetValueFromParams(raw)
		if originRow.After != nil {
			after = tools.SimpleMergeMap(*after, *originRow.After)
		}
		message := &model.OriginRow{
			Table:  fullName,
			Op:     model.OpType_Write,
			Before: &map[string]string{},
			After:  after,
			Full:   true,
		}
		//data, _ := json.Marshal(receiver)
		//tools.IssueLog("JSON [receiver]:%s", sql,  string(data))
		//data, _ = json.Marshal(message)
		//tools.IssueLog("JSON [message]:%s", sql,  string(data))
		// receiver <- message
		messages = append(messages, message)
	}
	return messages, nil
}
