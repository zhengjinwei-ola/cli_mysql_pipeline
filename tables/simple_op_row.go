package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type SimpleOpRowBase struct {
	Name        string
	Db          string
	DocField    string //文档依赖ID
	PkField     string //主键ID
	SqlTemplate string
	Fields      *[]model.EsField //如何格式数据
	Flow        model.FlowType   //操作类型
	FlowField   func(map[string]string) string
	Wrap        OpRowWrap
}

func (this SimpleOpRowBase) RetrieveRows(field string, val interface{}) []*model.OriginRow {

	var fullName = fmt.Sprintf("%s.%s", this.Db, this.Name)

	sql := fmt.Sprintf(
		"select * from %s where %s = ?",
		this.Name,
		field,
	)

	if len(this.SqlTemplate) > 0 {
		sql = this.SqlTemplate
	}

	res := []orm.Params{}
	_, err := model.Db.Raw(
		sql,
		val,
	).Values(&res)
	if err != nil {
		return nil
	}

	messages := make([]*model.OriginRow, 0, len(res))
	for _, raw := range res {
		message := &model.OriginRow{
			Table:  fullName,
			Op:     model.OpType_Write,
			Before: &map[string]string{},
			After:  model.GetValueFromParams(raw),
			Full:   true,
		}

		messages = append(messages, message)
	}
	return messages
}

func (this SimpleOpRowBase) RetrieveAll() []*model.OriginRow {
	if len(this.PkField) == 0 {
		return nil
	}
	//获取主键的最大值和最小值
	profile := model.TableProfile{}

	sql := fmt.Sprintf(
		"select max(%s) as max, min(%s) as min from %s",
		this.PkField,
		this.PkField,
		this.Name,
	)

	err := model.Db.Raw(sql).QueryRow(&profile)
	if err != nil {
		panic(err)
	}

	if profile.Max <= 0 {
		return nil
	}
	var min = profile.Min
	var max = profile.Max + 1

	var fullName = fmt.Sprintf("%s.%s", this.Db, this.Name)

	//根据主键范围，循环查询数据
	sql = fmt.Sprintf(
		"select * from %s where %s >= ? and %s < ?",
		this.Name,
		this.PkField,
		this.PkField,
	)

	if len(this.SqlTemplate) > 0 {
		sql = this.SqlTemplate
	}

	fmt.Printf("%s: max %d, min %d, max-min %d\n", fullName, max, min, max-min)

	res := []orm.Params{}
	_, err = model.Db.Raw(
		sql,
		min,
		max,
	).Values(&res)
	if err != nil {
		panic(err)
	}

	messages := make([]*model.OriginRow, 0, len(res))
	for _, raw := range res {
		message := &model.OriginRow{
			Table:  fullName,
			Op:     model.OpType_Write,
			Before: &map[string]string{},
			After:  model.GetValueFromParams(raw),
			Full:   true,
		}

		messages = append(messages, message)
	}
	return messages
}
