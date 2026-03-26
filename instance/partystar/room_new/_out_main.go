package main

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"fmt"
)

func main() {
	conf.ParseConfig("es-oversea-local")

	fmt.Println("start init db")
	model.InitDb()

	fmt.Println("start upgrade mapping")
	tables.UpgradeOverseaUserMapping()

	roomIndexer := index.RoomIndexer{
		Indexer: index.Indexer{
			Name:             "room_new",
			Mapping:          tables.MappingOverseaRoom,
			NumberOfShards:   2,
			NumberOfReplicas: 2,
		},
	}
	ops := []tables.SimpleOpRowBase{
		{
			Db:       "xianshi",
			Name:     "xs_chatroom",
			DocField: "rid",
			PkField:  "rid",
			Flow:     model.FlowType_Main, // FlowType_Main will create doc in index
			Fields: &[]model.EsField{
				{Name: "rid", Type: model.EsType_Number},
				{Name: "uid", Type: model.EsType_Number},
				{Name: "app_id", Type: model.EsType_Number},
				{Name: "name", Type: model.EsType_Text},
				{Name: "icon", Type: model.EsType_Text},
				{Name: "icon", Type: model.EsType_Text, Target: "room_icon"},
				{Name: "password", Type: model.EsType_TureOrFalse},
				{Name: "prefix", Type: model.EsType_Text},
				{Name: "property", Type: model.EsType_Text},
				{Name: "type", Type: model.EsType_Text},
				{Name: "types", Type: model.EsType_Text},
				{Name: "language", Type: model.EsType_Text},
				{Name: "game", Type: model.EsType_Text},
				{Name: "area", Type: model.EsType_Text},
				{Name: "dateline", Type: model.EsType_Number},
				{Name: "deleted", Type: model.EsType_Number},
				{Name: "weight", Type: model.EsType_Number},
				{Name: "room_factory_type", Type: model.EsType_Text, Target: "factory_type"},
				{Name: "sex", Type: model.EsType_Number, Target: "room_sex"},
				{Name: "paier", Type: model.EsType_Number},
				{Name: "settlement_channel", Type: model.EsType_Text},                         // 结算频道
				{Name: "fixed_tag_id", Type: model.EsType_Number, Target: "fixed_tag_id_new"}, // 外显标签
			},
			Wrap: tables.OpRowRoom{},
		},
		{
			Db:       "xianshi",
			Name:     "xs_chatroom_view",
			DocField: "rid",
			PkField:  "rid",
			Flow:     model.FlowType_Join, // FlowType_Main will update doc in index by given rid
			Fields: &[]model.EsField{
				{Name: "num_boy", Type: model.EsType_Number, Target: "boy_num"},
				{Name: "num_girl", Type: model.EsType_Number, Target: "girl_num"},
				{Name: "boss_uid", Type: model.EsType_Number},
				{Name: "mic_num", Type: model.EsType_Number},
				{Name: "hot", Type: model.EsType_Number, Target: "chatroom_view_hot_num"},
			},
			Wrap: tables.OpRowRoomView{},
		},
		{
			Db:       "xianshi",
			Name:     "xs_chatroom_config",
			DocField: "rid",
			PkField:  "rid",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "room_available_seat", Type: model.EsType_Number, Target: "room_available_seat"},
				{Name: "room_total_seat", Type: model.EsType_Number, Target: "room_total_seat"},
			},
			Wrap: tables.OpRowRoomConfigForRoom{},
		},
		{
			Db:       "xianshi",
			Name:     "es_update_lock_end",
			DocField: "rid",
			PkField:  "",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "rid", Type: model.EsType_Number},
				{Name: "lock_end", Type: model.EsType_Number},
			},
			Wrap: tables.OpRowRoomUpLockEnd{},
		},
		{
			Db:       "xianshi",
			Name:     "es_online_num",
			DocField: "rid",
			PkField:  "",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "rid", Type: model.EsType_Number},
				{Name: "online_num", Type: model.EsType_Number},
			},
			//可以直接使用基类的
			Wrap: tables.OpRowRoomUpLockEnd{},
		},
		{
			Db:       "xianshi",
			Name:     "es_room_real",
			DocField: "rid",
			PkField:  "",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "rid", Type: model.EsType_Number},
				{Name: "real", Type: model.EsType_Number},
			},
			//可以直接使用基类的
			Wrap: tables.OpRowTable{},
		},
		{
			Db:       "xianshi",
			Name:     "es_gaming_state",
			DocField: "rid",
			PkField:  "",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "gaming_name", Type: model.EsType_Text, Target: "gaming_name"},
				{Name: "min_players", Type: model.EsType_Number, Target: "min_players"},
				{Name: "gaming_state", Type: model.EsType_Number, Target: "gaming_state"},
			},
			Wrap: tables.OpRowTable{},
		},
		{
			Db:       "xianshi",
			Name:     "xs_chatroom_unity_gaming",
			DocField: "rid",
			PkField:  "rid",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "status", Type: model.EsType_Number, Target: "game_state_v2"},
				{Name: "dateline", Type: model.EsType_Number, Target: "game_state_create"},
				{Name: "updateline", Type: model.EsType_Number, Target: "game_state_update"},
				{Name: "version", Type: model.EsType_Number, Target: "game_version"},
			},
			//可以直接使用基类的
			Wrap: tables.OpRowTable{},
		},
		{
			Db:       "config",
			Name:     "bbc_home_page_tab_room",
			DocField: "rid",
			PkField:  "",
			Flow:     model.FlowType_Join,
			Fields: &[]model.EsField{
				{Name: "sort", Type: model.EsType_Number, Target: "pin"},
			},
			Wrap: tables.OpRowTable{},
		},
	}

	/*
		uncomment to generate room_new index in ES based on tables.MappingOverseaRoom if index is not already in ES
	*/
	// roomIndexer.DropIndex()
	// roomIndexer.CreateIndex()
	// roomIndexer.UpdateMapping()

	fmt.Println("num of tables", len(ops))

	for _, op := range ops {
		roomIndexer.UpdateRoomEs(105697419, op)
	}
}
