package room_new

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
)

func GetIndex() *index.Index {
	return &index.Index{
		Name:             "room_new",
		Receiver:         make(chan *model.OriginRow, 1),
		Mapping:          tables.MappingOverseaRoom,
		NumberOfShards:   2,
		NumberOfReplicas: 2,
		Ops: []tables.OpRowBase{
			{
				Db:        "xianshi",
				Name:      "xs_chatroom",
				IndexName: "ChatroomMainIndex",
				DocField:  "rid",
				PkField:   "rid",
				Flow:      model.FlowType_Main,
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
					// room state, 0: public, 1:private, 2: friends, 3:fans
					{Name: "state", Type: model.EsType_Number},
					{Name: "paier", Type: model.EsType_Number},
					{Name: "settlement_channel", Type: model.EsType_Text},                         // 结算频道
					{Name: "fixed_tag_id", Type: model.EsType_Number, Target: "fixed_tag_id_new"}, // 外显标签
					{Name: "tags", Type: model.EsType_Text},
				},
				Wrap:      tables.OpRowRoom{},
				UseUpsert: true,
			},
			// due to migration is very heavy, hence migration will use this to reduce the DB loads
			{
				Db:                 "xianshi",
				Name:               "es_xs_chatroom",
				IndexName:          "ChatroomStateIndex",
				DocField:           "rid",
				PkField:            "rid",
				SqlTemplate:        "SELECT rid, state FROM xs_chatroom WHERE rid >= ? and rid < ?",
				BatchSqlTemplate:   "SELECT rid, state FROM xs_chatroom WHERE rid IN (?)",
				MinMaxSqlTemplate:  "SELECT MAX(rid) as max, MIN(rid) as min FROM xs_chatroom",
				GappingSqlTemplate: "SELECT rid as min FROM xs_chatroom WHERE rid >= ? order by rid asc limit 1",
				Flow:               model.FlowType_Join,
				Fields: &[]model.EsField{
					// room state, 0: public, 1:private, 2: friends, 3:fans
					{Name: "state", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_chatroom_view",
				DocField: "rid",
				PkField:  "rid",
				Flow:     model.FlowType_Join,
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
				Fields:   &[]model.EsField{},
				Wrap:     tables.OpRowRoomConfigForRoom{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_chatroom_config",
				DocField: "rid",
				PkField:  "rid",
				Flow:     model.FlowType_Join,
				Fields:   &[]model.EsField{},
				Wrap:     tables.OpRowRoomConfigForRoomAge{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_chatroom_config",
				DocField: "rid",
				PkField:  "rid",
				Flow:     model.FlowType_Join,
				Fields:   &[]model.EsField{},
				Wrap:     tables.OpRowRoomConfigForScript{},
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
					{Name: "online_robot_num", Type: model.EsType_Number},
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
			// 每轮的房间流水
			{
				Db:       "xianshi",
				Name:     "es_room_round_real",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "round_real", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "es_room_round_real",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "round_real", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_bigarea",
				DocField: "rid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "bigarea_id", Type: model.EsType_Number},
				},
				Wrap:      tables.NewOpRowRoomBigArea("xs.bigarea.room"),
				UseUpsert: true,
			},
			{
				Db:        "xianshi",
				Name:      "es_room_gaming_level",
				IndexName: "EsRoomGamingLevel",
				DocField:  "rid",
				PkField:   "",
				Flow:      model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "room_gaming_level", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:        "xianshi",
				Name:      "es_room_block",
				IndexName: "EsRoomBlock",
				DocField:  "rid",
				PkField:   "",
				Flow:      model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "blocked", Type: model.EsType_TureOrFalse},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 玩法游戏 - 征招状态
			{
				Db:               "xianshi",
				Name:             "xs_user_match",
				IndexName:        "ChatRoomUserMatch",
				DocField:         "rid",
				PkField:          "rid",
				SqlTemplate:      "SELECT eee.rid, eee.uid AS recruit_uid, eee.top_category, (case when eee.status = 'wait' then 0 else 1 end) AS recruit_status, eee.version recruit_version from xs_user_match eee where eee.rid >= ? and eee.rid < ?",
				BatchSqlTemplate: "SELECT eee.rid, eee.uid AS recruit_uid, eee.top_category, (case when eee.status = 'wait' then 0 else 1 end) AS recruit_status, eee.version recruit_version from xs_user_match eee where eee.rid IN (?) ",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "recruit_uid", Type: model.EsType_Number},
					{Name: "top_category", Type: model.EsType_Number},
					{Name: "recruit_status", Type: model.EsType_Number},
					{Name: "recruit_version", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 玩法游戏 - 其他数据状态类型相关
			{
				Db:               "xianshi",
				Name:             "xs_chatroom",
				IndexName:        "RoomGameState",
				DocField:         "rid",
				PkField:          "rid",
				SqlTemplate:      "SELECT rid, -1, -1 AS call_state from xs_chatroom where rid >= ? and rid < ?",
				BatchSqlTemplate: "SELECT rid, -1, -1 AS call_state from xs_chatroom where rid IN (?)",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "call_state", Type: model.EsType_Number},
					{Name: "call_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:                 "xianshi",
				Name:               "es_xs_chatroom_config",
				IndexName:          "RoomAvailableSeatIndex",
				DocField:           "rid",
				PkField:            "rid",
				SqlTemplate:        "SELECT rid, SUM(CASE when uid = 0 and `lock` = 0 and forbidden = 0 then 1 else 0 end) as `room_available_seat`, SUM(CASE when uid <> 0 then 1 else 0 end) as `room_occupied_seat`, count(position) as `room_total_seat` FROM xs_chatroom_config WHERE rid >= ? and rid < ? group by rid",
				BatchSqlTemplate:   "SELECT rid, SUM(CASE when uid = 0 and `lock` = 0 and forbidden = 0 then 1 else 0 end) as `room_available_seat`, SUM(CASE when uid <> 0 then 1 else 0 end) as `room_occupied_seat`, count(position) as `room_total_seat` FROM xs_chatroom_config WHERE rid IN (?) group by rid",
				MinMaxSqlTemplate:  "SELECT MAX(ee.rid) as max, MIN(ee.rid) as min FROM xs_chatroom_config AS ee",
				GappingSqlTemplate: "SELECT ee.rid as min FROM xs_chatroom_config AS ee WHERE ee.rid >= ? order by rid asc limit 1",
				Flow:               model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "room_available_seat", Type: model.EsType_Number},
					{Name: "room_occupied_seat", Type: model.EsType_Number},
					{Name: "room_total_seat", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 你画我猜
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_guess",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT eee.rid, 'guess' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, 4 AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_guess eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'guess' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT eee.rid, 'guess' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, 4 AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_guess eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'guess' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// 剧本杀
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_juben",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT eee.rid, 'juben' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, eee.people_count AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_juben eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'puzzle' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT eee.rid, 'juben' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, eee.people_count AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_juben eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'puzzle' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// 海龟汤
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_puzzle",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT eee.rid, 'puzzle' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, eee.people_count AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_puzzle eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'puzzle' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT eee.rid, 'puzzle' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, eee.people_count AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_puzzle eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'puzzle' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// 谁是卧底
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_under",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT eee.rid, 'under' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, 5 AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_under eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'under' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT eee.rid, 'under' AS gaming_name, (case when eee.state = 'wait' then 0 else 1 end) AS gaming_state, 5 AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_under eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'under' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// 狼人杀 6，9，12
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_wolf",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT  eee.rid, 'wolf' AS gaming_name, (CASE WHEN eee.state = 'wait' THEN 0 ELSE 1 END) AS gaming_state, (CASE  WHEN eee.player_num_type IN (2,5) THEN 6 WHEN eee.player_num_type IN (3) THEN 12 ELSE 9 END) AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_wolf eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'wolf' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT  eee.rid, 'wolf' AS gaming_name, (CASE WHEN eee.state = 'wait' THEN 0 ELSE 1 END) AS gaming_state, (CASE  WHEN eee.player_num_type IN (2,5) THEN 6 WHEN eee.player_num_type IN (3) THEN 12 ELSE 9 END) AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_wolf eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'wolf' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// C位抢唱
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_grabmic",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT  eee.rid, 'grabmic' AS gaming_name, (CASE WHEN eee.state = 'wait' THEN 0 ELSE 1 END) AS gaming_state,  2 AS min_players, eee.next_state AS next_state, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_grabmic eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'grabmic' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT  eee.rid, 'grabmic' AS gaming_name, (CASE WHEN eee.state = 'wait' THEN 0 ELSE 1 END) AS gaming_state,  2 AS min_players, eee.next_state AS next_state, fff.position_num AS total_players " +
					" FROM xs_chatroom_extend_grabmic eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'grabmic' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// 残酷2选1
			{
				Db:        "xianshi",
				Name:      "xs_majority_extend",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT  eee.rid, 'majority' AS gaming_name, (CASE WHEN eee.stage = 0 THEN 0 ELSE 1 END) AS gaming_state, (CASE WHEN eee.mod = 1 THEN 2 ELSE 3 END) AS min_players, fff.position_num AS total_players " +
					" FROM xs_majority_extend eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'majority' WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT  eee.rid, 'majority' AS gaming_name, (CASE WHEN eee.stage = 0 THEN 0 ELSE 1 END) AS gaming_state, (CASE WHEN eee.mod = 1 THEN 2 ELSE 3 END) AS min_players, fff.position_num AS total_players " +
					" FROM xs_majority_extend eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'majority' WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// Unity 游戏类
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_unity_gaming",
				IndexName: "EsGameGamingState",
				DocField:  "rid",
				PkField:   "rid",
				SqlTemplate: "SELECT eee.rid, LOWER(eee.gaming) AS gaming_name, eee.status AS gaming_state, 2 AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_unity_gaming eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'unity' AND fff.room_game = eee.gaming WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate: "SELECT eee.rid, LOWER(eee.gaming) AS gaming_name, eee.status AS gaming_state, 2 AS min_players, fff.position_num AS total_players " +
					" FROM xs_chatroom_unity_gaming eee JOIN xs_chatroom_module_config fff ON fff.room_type = 'unity' AND fff.room_game = eee.gaming WHERE eee.rid IN (?) ",
				Flow: model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "gaming_name", Type: model.EsType_Text},
					{Name: "gaming_state", Type: model.EsType_Number},
					{Name: "min_players", Type: model.EsType_Number},
					{Name: "total_players", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// Unity 组自己的
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
			// 狼人杀历史数据
			{
				Db:                 "xianshi",
				Name:               "xs_chatroom_extend_wolf_history",
				IndexName:          "RoomGameLevel",
				DocField:           "rid",
				PkField:            "rid",
				SqlTemplate:        "select eee.rid, eee.uid, ff.level as wolf_level from xs_chatroom eee join xs_chatroom_extend_wolf_history ff ON ff.uid = eee.uid WHERE eee.rid >= ? and eee.rid < ? ",
				BatchSqlTemplate:   "select eee.rid, eee.uid, ff.level as wolf_level from xs_chatroom eee join xs_chatroom_extend_wolf_history ff ON ff.uid = eee.uid WHERE eee.rid IN (?) ",
				MinMaxSqlTemplate:  "select min(eee.rid) AS min, max(eee.rid) AS max from xs_chatroom eee join xs_chatroom_extend_wolf_history ff ON ff.uid = eee.uid GROUP BY eee.rid",
				GappingSqlTemplate: "select min(eee.rid) AS min from xs_chatroom eee join xs_chatroom_extend_wolf_history ff ON ff.uid = eee.uid GROUP BY eee.rid ORDER BY eee.rid asc limit 1",
				Flow:               model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "wolf_level", Type: model.EsType_Number, Target: "wolf_level"},
				},
				Wrap: tables.OpRowTable{},
			},
			// TODO: to be deleted
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
			// 房间连麦状态
			{
				Db:       "xianshi",
				Name:     "es_room_link_mic",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "link_mic_status", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "es_room_pk_state",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "pk_state", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "es_room_red_packet_post_time",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "red_packet_post_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "es_room_family_id",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "family_id", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房间团战PK状态
			{
				Db:       "xianshi",
				Name:     "es_room_team_pk_state",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "team_pk_state", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// feed 统计
			// 房间房主关注数量
			{
				Db:       "xianshi",
				Name:     "es_room_follow_num",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "follow_num", Type: model.EsType_Float},
					{Name: "follow_num_dateline", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房间评论数量
			{
				Db:       "xianshi",
				Name:     "es_room_comment_num",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "comment_num", Type: model.EsType_Float},
					{Name: "comment_num_dateline", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房间连麦数量
			{
				Db:       "xianshi",
				Name:     "es_room_link_mic_num",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "link_mic_num", Type: model.EsType_Float},
					{Name: "link_mic_num_dateline", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房间流水
			{
				Db:       "xianshi",
				Name:     "es_room_money",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "money", Type: model.EsType_Float},
					{Name: "money_dateline", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 更新查询条件 （兜底）
			{
				Db:       "xianshi",
				Name:     "es_room_query_condition",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "app_id", Type: model.EsType_Number},
					{Name: "uid", Type: model.EsType_Number},
					{Name: "language", Type: model.EsType_Text},
					{Name: "room_sex", Type: model.EsType_Number},
					{Name: "property", Type: model.EsType_Text},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 更新查询条件 （兜底）
			{
				Db:       "xianshi",
				Name:     "es_room_boom_rocket",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "boom_rocket_lv", Type: model.EsType_Number},
					{Name: "boom_rocket_failure_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 爆火箭加分
			{
				Db:       "xianshi",
				Name:     "es_room_game_money_type",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "game_money_type", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 使用置顶卡
			{
				Db:       "xianshi",
				Name:     "es_room_top_card",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "hot_num_score", Type: model.EsType_Number},
					{Name: "hot_num_score_expired_time", Type: model.EsType_Number},
					{Name: "hot_num_score_update_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 使用热门礼物
			{
				Db:       "xianshi",
				Name:     "es_feed_gift",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "feed_gift_hot_num_score", Type: model.EsType_Number},
					{Name: "feed_gift_hot_num_score_expired_time", Type: model.EsType_Number},
					{Name: "feed_gift_hot_num_score_update_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 多人pk
			{
				Db:       "xianshi",
				Name:     "es_room_multi_pk_state",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "multi_pk_status", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 多人连麦
			{
				Db:       "xianshi",
				Name:     "es_room_multi_link_mic_status",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "multi_link_mic_status", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// pk设置开关
			{
				Db:       "xianshi",
				Name:     "es_room_video_pk_settings",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "video_pk_settings", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房主当前开播版本号
			{
				Db:       "xianshi",
				Name:     "es_room_app_version",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "app_version", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房主当前邀请状态 1 空闲 0 已邀请
			{
				Db:       "xianshi",
				Name:     "es_room_multi_anchor_invite_status",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "multi_anchor_invite_status", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},

			{
				Db:       "xianshi",
				Name:     "es_room_multi_anchor_level",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "multi_anchor_level", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},

			{
				Db:       "xianshi",
				Name:     "es_room_is_ordinary",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "is_ordinary", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房间置底配置
			{
				Db:       "xianshi",
				Name:     "es_room_bottom_config",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "room_bottom_expired_time", Type: model.EsType_Number},
					{Name: "room_bottom_update_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 房间置底配置
			{
				Db:       "xianshi",
				Name:     "es_room_bottom_expired_time",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "room_bottom_expired_time", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
			// 语音房跨房pk状态
			{
				Db:       "xianshi",
				Name:     "es_room_chatroom_pk_status",
				DocField: "rid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number},
					{Name: "chatroom_pk_status", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowTable{},
			},
		},
	}
}
