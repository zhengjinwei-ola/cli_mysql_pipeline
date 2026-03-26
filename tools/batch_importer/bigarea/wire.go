// +build wireinject

package bigarea

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/common"
	"github.com/olachat/banban_server/common/nsq"
	"github.com/google/wire"
)

func InitNsqHandler(
	pkColName string,
	tables []string,
	idx *index.Index,
) *common.BatchImportHandler {
	wire.Build(
		common.NewHandler, common.NewBatchImportUseCase, newNsqAdapter,
		wire.Bind(new(common.UseCase), new(*common.BatchImportUseCase)),
		wire.Bind(new(common.TableStorage), new(*nsqAdapter)),
		wire.Bind(new(common.MsgAdapter), new(*nsqAdapter)),
	)
	return &common.BatchImportHandler{}
}

func InitNsqConsumer(
	topic nsq.Topic,
	channle nsq.Channel,
	h nsq.NsqHandleMessage,
) *nsq.NsqConsumer {
	wire.Build(nsq.NewNsqConsumer)
	return &nsq.NsqConsumer{}
}
