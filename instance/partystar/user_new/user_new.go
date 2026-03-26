package user_new

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
)

func GetIndex() *index.Index {
	return &index.Index{
		Name:             "user_new",
		Receiver:         make(chan *model.OriginRow, 1),
		Mapping:          tables.MappingOverseaUser,
		NumberOfShards:   3,
		NumberOfReplicas: 2,
		Ops: []tables.OpRowBase{
			{
				Db:          "xianshi",
				Name:        "xs_user_profile",
				DocField:    "uid",
				PkField:     "uid",
				SqlTemplate: "select p.*, e.interests from xs_user_profile as p left join xs_user_exposure as e on p.uid = e.uid where p.uid >= ? and p.uid < ?",
				Flow:        model.FlowType_Main,
				Fields: &[]model.EsField{
					{Name: "uid", Type: model.EsType_Number},
					{Name: "app_id", Type: model.EsType_Number},
					{Name: "name", Type: model.EsType_Text},
					{Name: "icon", Type: model.EsType_Text},
					{Name: "sign", Type: model.EsType_Text},
					{Name: "city", Type: model.EsType_Text},
					{Name: "position", Type: model.EsType_Text},
					{Name: "birthday", Type: model.EsType_Number},
					{Name: "job", Type: model.EsType_Number},
					{Name: "sex", Type: model.EsType_Number},
					{Name: "role", Type: model.EsType_Number},
					{Name: "god_category", Type: model.EsType_SetNumber},
					{Name: "god_num", Type: model.EsType_Number},
					{Name: "god_month_num", Type: model.EsType_Number},
					{Name: "god_week_num", Type: model.EsType_Number},
					{Name: "god_day_num", Type: model.EsType_Number},
					{Name: "god_now_num", Type: model.EsType_Number},
					{Name: "god_dateline", Type: model.EsType_Number},
					{Name: "god_default_id", Type: model.EsType_Number},
					{Name: "god_default_cid", Type: model.EsType_Number},
					{Name: "online_status", Type: model.EsType_Number},
					{Name: "online_dateline", Type: model.EsType_Number},
					{Name: "city_code", Type: model.EsType_Number},
					{Name: "dateline", Type: model.EsType_Number},
					{Name: "deleted", Type: model.EsType_Number},
					{Name: "pay_num", Type: model.EsType_Number},
					{Name: "pay_money", Type: model.EsType_Number},
					{Name: "pay_room_money", Type: model.EsType_Number},
					{Name: "service_score", Type: model.EsType_Number},
					{Name: "service_busy", Type: model.EsType_Number},
					{Name: "block_un_auther_message", Type: model.EsType_TureOrFalse},
					{Name: "service_pause", Type: model.EsType_TureOrFalse},
					{Name: "notice_order", Type: model.EsType_TureOrFalse},
					{Name: "notice_game", Type: model.EsType_TureOrFalse},
					{Name: "has_video", Type: model.EsType_TureOrFalse},
					{Name: "tag", Type: model.EsType_Number},
					{Name: "pay_receive_today", Type: model.EsType_Number},
					{Name: "title", Type: model.EsType_Number},
					{Name: "friend_state", Type: model.EsType_Number},
					{
						Name: "online_day",
						Type: model.EsType_Func,
						Func: func(origin map[string]string) interface{} {
							val, ok := origin["online_dateline"]
							if ok {
								sec, err := strconv.ParseInt(val, 10, 64)
								if err == nil {
									date := time.Unix(sec, 0).Format("20060102")
									ymd, _ := strconv.ParseInt(date, 10, 64)
									return ymd
								}
							}
							return 0
						},
					},
					{
						Name: "geo",
						Type: model.EsType_Func,
						Func: func(origin map[string]string) interface{} {
							geo := model.EsGeo{}
							longitude, ok1 := origin["longitude"]
							latitude, ok2 := origin["latitude"]
							if ok1 && ok2 {
								geo.Lon, _ = strconv.ParseFloat(longitude, 64)
								geo.Lat, _ = strconv.ParseFloat(latitude, 64)
							}
							return geo
						},
					},
				},
				Wrap:      tables.OpRowUserProfile{},
				UseUpsert: true,
			},
			{
				Db:       "xianshi",
				Name:     "xs_chatroom_config",
				DocField: "uid",
				PkField:  "id",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "rid", Type: model.EsType_Number, Target: "room_rid"},
					{Name: "position", Type: model.EsType_Number, Target: "room_position"},
				},
				Wrap: tables.OpRowOverseaRoomConfigForUser{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_mentor_exp",
				DocField: "uid",
				PkField:  "id",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "level", Type: model.EsType_Number, Target: "mentor_level"},
					{Name: "left_disciple_num", Type: model.EsType_Number, Target: "left_disciple_num"},
				},
				Wrap: tables.OpRowMentorExp{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_interest_tags",
				DocField: "uid",
				PkField:  "",
				Flow:     model.FlowType_Append,
				FlowField: func(origin map[string]string) string {
					return "interests"
				},
				Wrap: tables.OpRowUserInterestTags{},
			},
			{
				Db:       "xianshi",
				Name:     "es_exposure",
				DocField: "uid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "uid", Type: model.EsType_Number},
					{Name: "num", Type: model.EsType_Number},
					{Name: "type", Type: model.EsType_Text},
				},
				Wrap: tables.OpRowEsExposure{},
			},
			{
				Db:        "xianshi",
				Name:      "es_viewing_page",
				IndexName: "EsUserViewingPage",
				DocField:  "uid",
				PkField:   "",
				Flow:      model.FlowType_Join,
				Fields: &[]model.EsField{
					// 1 = tier 1 page, 2 = tier 2 page and so on
					{Name: "viewing_page", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_settings",
				DocField: "uid",
				PkField:  "uid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "nearby_invisible", Type: model.EsType_TureOrFalse},
					{Name: "language", Type: model.EsType_Text},
				},
				Wrap: tables.OpRowUserSettings{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_country",
				DocField: "uid",
				PkField:  "uid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "country", Type: model.EsType_Text},
					{Name: "latest_country_code", Type: model.EsType_Text, Target: "country_code"},
				},
				Wrap: tables.OpRowUserSettings{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_god",
				DocField: "uid",
				PkField:  "id",
				Flow:     model.FlowType_Append,
				FlowField: func(origin map[string]string) string {
					return fmt.Sprintf("skill_%s", origin["cid"])
				},
				Wrap: tables.OpRowUserGod{},
				Fields: &[]model.EsField{
					{Name: "id", Type: model.EsType_Number},
					{Name: "cid", Type: model.EsType_Number},
					{Name: "level", Type: model.EsType_Number},
					{Name: "disabled", Type: model.EsType_TureOrFalse},
					{Name: "price", Type: model.EsType_Number},
					{Name: "discount", Type: model.EsType_Number},
					{Name: "tags", Type: model.EsType_SetNumber},
					{Name: "is_default", Type: model.EsType_TureOrFalse},
					{Name: "num", Type: model.EsType_Number},
					{Name: "now_num", Type: model.EsType_Number},
					{Name: "day_num", Type: model.EsType_Number},
					{Name: "week_num", Type: model.EsType_Number},
					{Name: "month_num", Type: model.EsType_Number},
					{Name: "day_credit", Type: model.EsType_Number},
					{Name: "week_credit", Type: model.EsType_Number},
					{Name: "now_credit", Type: model.EsType_Number},
					{Name: "closed", Type: model.EsType_TureOrFalse},
					{Name: "cover", Type: model.EsType_Text},
					{Name: "sign", Type: model.EsType_Text},
					{Name: "description", Type: model.EsType_Text},
					{Name: "audio", Type: model.EsType_Text},
					{
						Name: "pay_price",
						Type: model.EsType_Func,
						Func: func(origin map[string]string) interface{} {
							a, ok1 := origin["discount"]
							b, ok2 := origin["price"]
							if ok1 && ok2 {
								discount, _ := strconv.ParseInt(a, 10, 64)
								price, _ := strconv.ParseInt(b, 10, 64)
								if discount == 0 {
									discount = 10
								}
								priceFloat := float64(discount) * float64(price) / 10 / 100
								priceInt := int64(math.Round(priceFloat * 100))
								if priceInt < 100 {
									priceInt = 100
								}
								return priceInt
							}
							return 0
						},
					},
				},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_exposure",
				DocField: "uid",
				PkField:  "uid",
				Flow:     model.FlowType_Append,
				Fields: &[]model.EsField{
					{Name: "exposure_new", Type: model.EsType_Number},
					{Name: "exposure_old", Type: model.EsType_Number},
					{Name: "follow_new", Type: model.EsType_Number},
					{Name: "income_yesterday", Type: model.EsType_Number},
					{Name: "income_lastweek", Type: model.EsType_Number},
					{Name: "income_total", Type: model.EsType_Number},
					{Name: "is_peipei", Type: model.EsType_TureOrFalse},
					{Name: "follows", Type: model.EsType_SetNumber},
					{Name: "interests", Type: model.EsType_SetNumber},
					{Name: "friends_num", Type: model.EsType_Number},
					{Name: "clicked_uids", Type: model.EsType_SetNumber},
				},
				FlowField: func(origin map[string]string) string {
					return "exposure"
				},
				Wrap: tables.OpRowUserExposure{},
			},
			{
				Db:       "xianshi",
				Name:     "es_popularity",
				DocField: "uid",
				PkField:  "",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "uid", Type: model.EsType_Number},
					{Name: "popularity", Type: model.EsType_Number},
				},
				//可以直接使用基类的
				Wrap: tables.OpRowRoomUpLockEnd{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_pt_meet_settings",
				DocField: "uid",
				PkField:  "uid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "sex", Type: model.EsType_Number, Target: "meet_setting_sex"},
				},
				Wrap: tables.OpRowUserSettings{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_profile_extend",
				DocField: "uid",
				PkField:  "uid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "personal_tag", Type: model.EsType_SetNumber},
					{Name: "interest_tag", Type: model.EsType_SetNumber},
				},
				Wrap: tables.OpRowUserSettings{},
			},
			{
				Db:       "xianshi",
				Name:     "xs_user_bigarea",
				DocField: "uid",
				PkField:  "uid",
				Flow:     model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "bigarea_id", Type: model.EsType_Number},
				},
				Wrap:      tables.OpRowTable{},
				UseUpsert: true,
			},
			// 用户本身的游戏设定
			{
				Db:               "xianshi",
				Name:             "xs_user_extend_game_settings",
				IndexName:        "UserGameSettings",
				SqlTemplate:      "SELECT ee.uid, ee.mute_sys_game_prompts, ee.mute_stranger_game_invite, ee.mute_sys_game_prompts_till, ee.mute_stranger_game_invite_till FROM xs_user_extend_game_settings AS ee WHERE ee.uid >= ? AND ee.uid < ?",
				BatchSqlTemplate: "SELECT ee.uid, ee.mute_sys_game_prompts, ee.mute_stranger_game_invite, ee.mute_sys_game_prompts_till, ee.mute_stranger_game_invite_till FROM xs_user_extend_game_settings AS ee WHERE ee.uid IN (?)",
				DocField:         "uid",
				PkField:          "uid",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "mute_sys_game_prompts", Type: model.EsType_TureOrFalse},
					{Name: "mute_stranger_game_invite", Type: model.EsType_TureOrFalse},
					{Name: "mute_sys_game_prompts_till", Type: model.EsType_Number},
					{Name: "mute_stranger_game_invite_till", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
			// 是否已经推过游戏弹窗过给用户 (日后也可用作其他弹窗的防骚扰)
			{
				Db:               "xianshi",
				Name:             "xs_sys_prompt_user",
				IndexName:        "SysPromptUser",
				SqlTemplate:      "SELECT ee.uid, ee.prompted_to, ee.expire_at FROM xs_sys_prompt_user AS ee WHERE ee.uid >= ? AND ee.uid < ?",
				BatchSqlTemplate: "SELECT ee.uid, ee.prompted_to, ee.expire_at FROM xs_sys_prompt_user AS ee WHERE ee.uid IN (?)",
				DocField:         "uid",
				PkField:          "uid",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "prompted_to", Type: model.EsType_TureOrFalse, Target: "sys_game_match_prompted_to"},
					{Name: "expire_at", Type: model.EsType_Number, Target: "sys_game_match_expire_at"},
				},
				Wrap: tables.OpRowTable{},
			},
			// 画猜历史数据
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_guess_history",
				IndexName: "UserRoomGameHistory",
				DocField:  "uid",
				PkField:   "uid",
				Flow:      model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "total", Type: model.EsType_Number, Target: "recent_guess_total"},
					{Name: "score", Type: model.EsType_Number, Target: "recent_guess_score"},
					{Name: "thumbs_up", Type: model.EsType_Number, Target: "recent_guess_thumbs_up"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_guess_last_played"},
				},
				Wrap: tables.OpRowTable{},
			},
			// 剧本杀历史数据
			{
				Db:               "xianshi",
				Name:             "xs_user_juben_history",
				IndexName:        "UserRoomGameHistory",
				DocField:         "uid",
				PkField:          "uid",
				SqlTemplate:      "SELECT ee.uid, sum(ee.score) AS score, sum(ee.game_score) AS game_score, max(ee.dateline) AS dateline FROM xs_user_juben_history AS ee WHERE ee.uid >= ? AND ee.uid < ? GROUP BY ee.uid",
				BatchSqlTemplate: "SELECT ee.uid, sum(ee.score) AS score, sum(ee.game_score) AS game_score, max(ee.dateline) AS dateline FROM xs_user_juben_history AS ee WHERE ee.uid IN (?) GROUP BY ee.uid",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "score", Type: model.EsType_Number, Target: "recent_juben_score"},
					{Name: "game_score", Type: model.EsType_Number, Target: "recent_juben_game_score"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_juben_last_played"},
				},
				Wrap: tables.OpRowTable{},
			},
			// 海龟汤历史数据
			{
				Db:               "xianshi",
				Name:             "xs_user_puzzle_history",
				IndexName:        "UserRoomGameHistory",
				DocField:         "uid",
				PkField:          "uid",
				SqlTemplate:      "SELECT ee.uid, sum(ee.uid) AS total, sum(ee.score) AS score, max(ee.dateline) AS dateline FROM xs_user_puzzle_history AS ee WHERE ee.uid >= ? AND ee.uid < ? GROUP BY ee.uid",
				BatchSqlTemplate: "SELECT ee.uid, sum(ee.uid) AS total, sum(ee.score) AS score, max(ee.dateline) AS dateline FROM xs_user_puzzle_history AS ee WHERE ee.uid IN (?) GROUP BY ee.uid",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "total", Type: model.EsType_Number, Target: "recent_puzzle_total"},
					{Name: "score", Type: model.EsType_Number, Target: "recent_puzzle_score"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_puzzle_last_played"},
				},
				Wrap: tables.OpRowTable{},
			},
			// 卧底历史数据
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_under_history",
				IndexName: "UserRoomGameHistory",
				DocField:  "uid",
				PkField:   "uid",
				Flow:      model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "total", Type: model.EsType_Number, Target: "recent_under_total"},
					{Name: "populace", Type: model.EsType_Number, Target: "recent_under_populace"},
					{Name: "under", Type: model.EsType_Number, Target: "recent_under_under"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_under_last_played"},
				},
				Wrap: tables.OpRowTable{},
			},
			// 狼人杀历史数据
			{
				Db:        "xianshi",
				Name:      "xs_chatroom_extend_wolf_history",
				IndexName: "UserRoomGameHistory",
				DocField:  "uid",
				PkField:   "uid",
				Flow:      model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "total", Type: model.EsType_Number, Target: "recent_wolf_total"},
					{Name: "win", Type: model.EsType_Number, Target: "recent_wolf_win"},
					{Name: "goodWin", Type: model.EsType_Number, Target: "recent_wolf_goodWin"},
					{Name: "level", Type: model.EsType_Number, Target: "recent_wolf_level"},
					{Name: "exp", Type: model.EsType_Number, Target: "recent_wolf_exp"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_wolf_last_played"},
				},
				Wrap: tables.OpRowTable{},
			},
			// Unity - Ludo 历史数据
			{
				Db:                 "xianshi",
				Name:               "es_xs_unity_ludo_record",
				IndexName:          "UserUnityGameLudoHistory",
				DocField:           "uid",
				PkField:            "uid",
				SqlTemplate:        "SELECT ee.uid, ee.game, ee.round, ee.champion, ee.score, ee.updateline, ee.dateline FROM xs_unity_user_record AS ee WHERE ee.game = 'ludo' AND ee.uid >= ? AND ee.uid < ? ",
				BatchSqlTemplate:   "SELECT ee.uid, ee.game, ee.round, ee.champion, ee.score, ee.updateline, ee.dateline AS dateline FROM xs_unity_user_record AS ee WHERE ee.game = 'ludo' AND ee.uid IN (?) ",
				MinMaxSqlTemplate:  "SELECT MAX(ee.uid) as max, MIN(ee.uid) as min FROM xs_unity_user_record AS ee WHERE ee.game = 'ludo' ",
				GappingSqlTemplate: "SELECT ee.uid as min FROM xs_unity_user_record AS ee WHERE ee.game = 'ludo' AND uid >= ? order by uid asc limit 1",
				Flow:               model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "round", Type: model.EsType_Number, Target: "recent_ludo_round"},
					{Name: "champion", Type: model.EsType_Number, Target: "recent_ludo_champion"},
					{Name: "score", Type: model.EsType_Number, Target: "recent_ludo_score"},
					{Name: "updateline", Type: model.EsType_Number, Target: "recent_ludo_last_played"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_ludo_dateline"},
				},
				Wrap: tables.OpRowTable{},
			},
			// Unity - Carrom 历史数据
			{
				Db:                 "xianshi",
				Name:               "es_xs_unity_carrom_record",
				IndexName:          "UserUnityGameCarromHistory",
				DocField:           "uid",
				PkField:            "uid",
				SqlTemplate:        "SELECT ee.uid, ee.game, ee.round, ee.champion, ee.score, ee.updateline, ee.dateline FROM xs_unity_user_record AS ee WHERE ee.game = 'carrom' AND ee.uid >= ? AND ee.uid < ? ",
				BatchSqlTemplate:   "SELECT ee.uid, ee.game, ee.round, ee.champion, ee.score, ee.updateline, ee.dateline AS dateline FROM xs_unity_user_record AS ee WHERE ee.game = 'carrom' AND ee.uid IN (?) ",
				MinMaxSqlTemplate:  "SELECT MAX(ee.uid) as max, MIN(ee.uid) as min FROM xs_unity_user_record AS ee WHERE ee.game = 'carrom' ",
				GappingSqlTemplate: "SELECT ee.uid as min FROM xs_unity_user_record AS ee WHERE ee.game = 'carrom' AND uid >= ? order by uid asc limit 1",
				Flow:               model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "round", Type: model.EsType_Number, Target: "recent_carrom_round"},
					{Name: "champion", Type: model.EsType_Number, Target: "recent_carrom_champion"},
					{Name: "score", Type: model.EsType_Number, Target: "recent_carrom_score"},
					{Name: "updateline", Type: model.EsType_Number, Target: "recent_carrom_last_played"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_carrom_dateline"},
				},
				Wrap: tables.OpRowTable{},
			},
			// Unity - Billiards 历史数据
			{
				Db:                 "xianshi",
				Name:               "es_xs_unity_billiards_record",
				IndexName:          "UserUnityGameBilliardsHistory",
				DocField:           "uid",
				PkField:            "uid",
				SqlTemplate:        "SELECT ee.uid, ee.game, ee.round, ee.champion, ee.score, ee.updateline, ee.dateline FROM xs_unity_user_record AS ee WHERE ee.game = 'billiards' AND ee.uid >= ? AND ee.uid < ? ",
				BatchSqlTemplate:   "SELECT ee.uid, ee.game, ee.round, ee.champion, ee.score, ee.updateline, ee.dateline AS dateline FROM xs_unity_user_record AS ee WHERE ee.game = 'billiards' AND ee.uid IN (?) ",
				MinMaxSqlTemplate:  "SELECT MAX(ee.uid) as max, MIN(ee.uid) as min FROM xs_unity_user_record AS ee WHERE ee.game = 'billiards' ",
				GappingSqlTemplate: "SELECT ee.uid as min FROM xs_unity_user_record AS ee WHERE ee.game = 'billiards' AND uid >= ? order by uid asc limit 1",
				Flow:               model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "round", Type: model.EsType_Number, Target: "recent_billiards_round"},
					{Name: "champion", Type: model.EsType_Number, Target: "recent_billiards_champion"},
					{Name: "score", Type: model.EsType_Number, Target: "recent_billiards_score"},
					{Name: "updateline", Type: model.EsType_Number, Target: "recent_billiards_last_played"},
					{Name: "dateline", Type: model.EsType_Number, Target: "recent_billiards_dateline"},
				},
				Wrap: tables.OpRowTable{},
			},
			// C位抢唱历史数据
			{
				Db:               "xianshi",
				Name:             "xs_grabmic_user_song",
				IndexName:        "UserRoomGameHistory",
				DocField:         "uid",
				PkField:          "uid",
				SqlTemplate:      "SELECT ee.uid, SUM(ee.count) AS count,  SUM(CASE WHEN ee.status = 0 THEN 1 ELSE 0 END) AS pending_confirm, SUM(CASE WHEN ee.status  = 1 THEN 1 ELSE 0 END) AS success, SUM(CASE WHEN ee.status  = 2 THEN 1 ELSE 0 END) AS fail, max(CASE WHEN ee.create_time > ee.update_time THEN ee.create_time ELSE ee.update_time END) AS last_update_time, max(ee.create_time) AS last_create_time FROM xs_grabmic_user_song  ee  WHERE ee.uid >= ? AND ee.uid < ? GROUP BY ee.uid",
				BatchSqlTemplate: "SELECT ee.uid, SUM(ee.count) AS count,  SUM(CASE WHEN ee.status = 0 THEN 1 ELSE 0 END) AS pending_confirm, SUM(CASE WHEN ee.status  = 1 THEN 1 ELSE 0 END) AS success, SUM(CASE WHEN ee.status  = 2 THEN 1 ELSE 0 END) AS fail, max(CASE WHEN ee.create_time > ee.update_time THEN ee.create_time ELSE ee.update_time END) AS last_update_time, max(ee.create_time) AS last_create_time FROM xs_grabmic_user_song  ee  WHERE ee.uid IN (?) GROUP BY ee.uid",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "count", Type: model.EsType_Number, Target: "recent_grabmic_total"},
					{Name: "pending_confirm", Type: model.EsType_Number, Target: "recent_grabmic_pending_confirm"},
					{Name: "success", Type: model.EsType_Number, Target: "recent_grabmic_success"},
					{Name: "fail", Type: model.EsType_Number, Target: "recent_grabmic_fail"},
					{Name: "last_update_time", Type: model.EsType_Number, Target: "recent_grabmic_last_played"},
					{Name: "last_create_time", Type: model.EsType_Number, Target: "recent_grabmic_last_join"},
				},
				Wrap: tables.OpRowTable{},
			},
			// 残酷2选1历史数据
			{
				Db:               "xianshi",
				Name:             "xs_majority_result",
				IndexName:        "UserRoomGameHistory",
				DocField:         "uid",
				PkField:          "uid",
				SqlTemplate:      "SELECT ee.uid, count(*) AS total,  SUM(ee.score) AS total_score, max(CASE WHEN ee.update_time > ee.create_time THEN ee.update_time ELSE ee.create_time END) AS last_update_time, max(ee.create_time) AS last_create_time FROM xs_majority_result ee WHERE ee.uid >= ? AND ee.uid < ? GROUP BY ee.uid",
				BatchSqlTemplate: "SELECT ee.uid, count(*) AS total,  SUM(ee.score) AS total_score, max(CASE WHEN ee.update_time > ee.create_time THEN ee.update_time ELSE ee.create_time END) AS last_update_time, max(ee.create_time) AS last_create_time FROM xs_majority_result ee WHERE ee.uid IN (?) GROUP BY ee.uid",
				Flow:             model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "total", Type: model.EsType_Number, Target: "recent_majority_total"},
					{Name: "total_score", Type: model.EsType_Number, Target: "recent_majority_total_score"},
					{Name: "last_update_time", Type: model.EsType_Number, Target: "recent_majority_last_played"},
					{Name: "last_create_time", Type: model.EsType_Number, Target: "recent_majority_last_join"},
				},
				Wrap: tables.OpRowTable{},
			},
		},
	}
}
