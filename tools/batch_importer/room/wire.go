// +build wireinject

package room

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/common"
	"github.com/olachat/banban_server/common/nsq"
	"github.com/astaxie/beego/orm"
	"github.com/google/wire"
)

func provideNsqConsumer(
	topic nsq.Topic,
	channel nsq.Channel,
	h *common.BatchImportHandler,
) *nsq.NsqConsumer {
	return nsq.NewNsqConsumer(topic, channel, h.HandleMsg)
}

func InitNsqConsumer(
	topic nsq.Topic,
	channel nsq.Channel,
	pkColName string,
	tables []string,
	sql SQL,
	idx *index.Index,
	db orm.Ormer,
) *nsq.NsqConsumer {
	wire.Build(
		provideNsqConsumer, common.NewHandler, common.NewBatchImportUseCase, newPHPAdapter, newMysqlAdapter,
		wire.Bind(new(common.UseCase), new(*common.BatchImportUseCase)),
		wire.Bind(new(common.TableStorage), new(*mysqlAdapter)),
		wire.Bind(new(common.MsgAdapter), new(*phpAdapter)),
	)
	return &nsq.NsqConsumer{}
}
