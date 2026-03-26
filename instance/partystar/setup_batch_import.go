package main

import (
	"strings"

	"github.com/kpango/glg"
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/bigarea"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/room"
)

var (
	roomTemplate = "select r.rid, r.sort from config.bbc_home_page_tab_room as r where r.rid in (%v)"
)

func setupBatchImport(roomIdx, userIdx *index.Index) {
	temp := strings.Split(conf.NsqAddr, ",")
	addrs := make([]string, 0)
	for _, address := range temp {
		addrs = append(addrs, strings.TrimSpace(address))
	}
	glg.Info("conf item nsq", conf.NsqAddr)

	roomConsumer := room.InitNsqConsumer(
		"xs.room.promote", "xs_chatroom", "rid",
		[]string{"config.bbc_home_page_tab_room"},
		room.SQL(roomTemplate), roomIdx, model.Db,
	)

	bigAreaUserHandler := bigarea.InitNsqHandler(
		"uid", []string{"xianshi.xs_user_bigarea"}, userIdx,
	)

	bigAreaRoomHandler := bigarea.InitNsqHandler(
		"rid", []string{"xianshi.xs_user_bigarea"}, roomIdx,
	)

	bigAreaConsumer := bigarea.InitNsqConsumer(
		"xs.property",
		"bigarea",
		bigarea.WrapBigAreaHandlers(bigAreaRoomHandler, bigAreaUserHandler),
	)

	_ = roomConsumer.Connect(addrs)
	_ = bigAreaConsumer.Connect(addrs)
}
