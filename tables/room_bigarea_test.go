package tables_test

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/app"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"github.com/olachat/banban_server/common/nsq"
	"github.com/olachat/banban_server/common/serialize"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	goNsq "github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqd"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	room    *index.Index
	control *app.App
	channel chan *model.OriginRow
	got     []interface{}
	want    []interface{}
)

func TestOpRowRoomBigArea_Format(t *testing.T) {
	t.Run("nsq publish", func(t *testing.T) {
		want = append(want, map[string]interface{}{"rid": "1111", "bigarea_id": "0"})
		want = append(want, map[string]interface{}{"rid": "2222", "bigarea_id": "0"})
		got = make([]interface{}, 0)
		setupDB()
		setupNSQ()
		setupApp()
		insertRecord()
		stopApp()

		time.Sleep(1 * time.Second)
		assert.ElementsMatch(t, want, got)
	})
}

func stopApp() {
	control.Stop()
}

func setupNSQ() {
	conf.NsqChannel = "test"
	conf.NsqAddr = "0.0.0.0:4150"
	opt := nsqd.NewOptions()
	nsqd, _ := nsqd.New(opt)
	go nsqd.Main()
	consumer := nsq.NewNsqConsumer("testtopic", "test", func(msg *goNsq.Message) error {
		v, _ := serialize.UnMarshal(msg.Body)
		got = append(got, v)
		return nil
	})
	consumer.Connect([]string{conf.NsqAddr})
}

func insertRecord() {
	channel <- &model.OriginRow{
		Table:  "xianshi.xs_user_bigarea",
		Op:     model.OpType_Write,
		Before: &map[string]string{},
		After: &map[string]string{
			"uid":        "1234567",
			"bigarea_id": "0",
		},
	}
}

func setupApp() {
	channel = make(chan *model.OriginRow, 2)
	room = &index.Index{
		Name:             "room_new",
		Receiver:         channel,
		Mapping:          tables.MappingOverseaRoom,
		NumberOfShards:   2,
		NumberOfReplicas: 2,
		Ops: []tables.OpRowBase{
			{
				Db:       "xianshi",
				Name:     "xs_user_bigarea",
				DocField: "rid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "bigarea_id", Type: model.EsType_Number},
				},
				Wrap: tables.NewOpRowRoomBigArea("testtopic"),
			},
		},
	}
	control = app.NewApp()
	control.Add(room)
	go control.Start("", "")
}

func setupDB() {
	engine := sqle.NewDefault()
	engine.AddDatabase(information_schema.NewInformationSchemaDatabase(engine.Catalog))

	db := createTestDatabase()
	engine.AddDatabase(db)

	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Auth:     auth.NewNativeSingle("root", "123456", auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}
	go s.Start()

	conf.MysqlAddr = "root:123456@tcp(127.0.0.1:3306)/xianshi?charset=utf8mb4&parseTime=true"
	conf.MysqlMasterAddr = "root:123456@tcp(127.0.0.1:3306)/xianshi?charset=utf8mb4&parseTime=true"
	model.InitDb()
}

func createTestDatabase() *memory.Database {
	const (
		dbName    = "xianshi"
		tableName = "xs_chatroom"
	)

	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "uid", Type: sql.Int32, Nullable: false, Source: tableName},
		{Name: "rid", Type: sql.Int32, Nullable: false, Source: tableName, PrimaryKey: true},
		{Name: "property", Type: sql.Text, Nullable: true, Source: tableName},
		{Name: "type", Type: sql.Text, Nullable: true, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()
	_ = table.Insert(ctx, sql.NewRow(1234567, 1111, "test", "test"))
	_ = table.Insert(ctx, sql.NewRow(1234567, 2222, "test", "test"))
	return db
}
