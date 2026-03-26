package user_god_new

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"fmt"
	"math"
	"strconv"
	"time"
)

func GetIndex() *index.Index {
	return &index.Index{
		Name:             "user_god_new",
		Receiver:         make(chan *model.OriginRow, 1),
		Mapping:          tables.MappingOverseaUser,
		NumberOfShards:   2,
		NumberOfReplicas: 3,
		Ops: []tables.OpRowBase{
			{
				Db:          "xianshi",
				Name:        "xs_user_profile",
				DocField:    "uid",
				PkField:     "uid",
				SqlTemplate: "select p.*, e.interests from xs_user_profile as p left join xs_user_exposure as e on p.uid = e.uid where p.uid >= ? and p.uid < ? and p.role >= 2",
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
				Wrap: tables.OpRowUserProfileForGod{},
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
				Db:          "xianshi",
				Name:        "xs_user_settings",
				SqlTemplate: "select e.* from xs_user_profile as p left join xs_user_settings as e on p.uid = e.uid where p.uid >= ? and p.uid < ? and p.role >= 2",
				DocField:    "uid",
				PkField:     "uid",
				Flow:        model.FlowType_Join,
				Fields: &[]model.EsField{
					{Name: "nearby_invisible", Type: model.EsType_TureOrFalse},
					{Name: "language", Type: model.EsType_Text},
				},
				Wrap: tables.OpRowUserSettings{},
			},
			{
				Db:          "xianshi",
				Name:        "xs_user_country",
				SqlTemplate: "select e.* from xs_user_profile as p left join xs_user_country as e on p.uid = e.uid where p.uid >= ? and p.uid < ? and p.role >= 2",
				DocField:    "uid",
				PkField:     "uid",
				Flow:        model.FlowType_Join,
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
				Db:          "xianshi",
				Name:        "xs_user_exposure",
				DocField:    "uid",
				PkField:     "uid",
				SqlTemplate: "select e.* from xs_user_profile as p left join xs_user_exposure as e on p.uid = e.uid where p.uid >= ? and p.uid < ? and p.role >= 2",
				Flow:        model.FlowType_Append,
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
		},
	}
}
