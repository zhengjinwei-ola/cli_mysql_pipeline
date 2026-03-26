package index

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"fmt"
)

type SimpleIndex struct {
	Ops    []tables.SimpleOpRowBase //哪些表使用哪个类去操作
	opsMap map[string][]tables.SimpleOpRowBase
	Indexer
}

func (this *SimpleIndex) Init() {
	this.opsMap = make(map[string][]tables.SimpleOpRowBase)
	for _, op := range this.Ops {
		name := fmt.Sprintf("%s.%s", op.Db, op.Name)
		if _, ok := this.opsMap[name]; !ok {
			this.opsMap[name] = make([]tables.SimpleOpRowBase, 0)
		}
		this.opsMap[name] = append(this.opsMap[name], op)
	}
}

func (this *SimpleIndex) StartIndexing(autoDropIndex bool) {
	//尝试创建index
	if autoDropIndex {
		fmt.Println("drop start", this.Name)
		this.DropIndex()
	}

	this.CreateIndex()
	this.UpdateMapping()
	for _, op := range this.Ops {
		messages := op.RetrieveAll()
		var opOutput []string
		for _, msg := range messages {
			output, err := this.ProcessOriginRow(msg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			opOutput = append(opOutput, output...)
		}
		fmt.Println(op.Name, "has", len(opOutput), "strings")
		bulkIngest(opOutput)
	}
}

func (this *SimpleIndex) ProcessOriginRow(row *model.OriginRow) ([]string, error) {
	name := row.Table
	ops, ok := this.opsMap[name]
	if !ok {
		return nil, fmt.Errorf("SimpleIndex %s Table not found %s", this.Name, name)
	}
	output := make([]string, 0, len(ops))
	for _, op := range ops {
		if row.Full {
			op.Wrap.InitHookBefore(row.After)
		}
		docId, doc, opType := op.Wrap.Format(
			*row.After,
			*row.Before,
			op.DocField,
			op.Fields,
			row.Op,
		)
		if docId == 0 || doc == nil {
			continue
		}

		if row.Full {
			op.Wrap.InitHook(docId, &doc)
		}
		message := &model.Message{
			Id:   docId,
			Data: &doc,
			Op:   opType,
			Flow: op.Flow,
			Pk:   op.PkField,
		}
		if op.Flow == model.FlowType_Append {
			message.FlowField = op.FlowField(*row.After)
		}

		//bytes, _ := json.Marshal(message.Data)
		//log.Println("doc", message.Id, len(string(bytes)))
		//根据flow和op写入数据
		var err error
		var o string
		switch message.Flow {
		case model.FlowType_Main:
			o, err = this.flowMain(message)
			break

		case model.FlowType_Join:
			o, err = this.flowJoin(message)
			break

		case model.FlowType_Append:
			o, err = this.flowAppend(message)
			break
		default:
			err = fmt.Errorf("unknown message flow")
		}
		if err == nil {
			output = append(output, o)
		}
	}
	return output, nil
}
