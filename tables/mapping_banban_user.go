package tables

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"fmt"
)

/*
	geo
	字符串形式以半角逗号分割，如 "lat,lon"
	对象形式显式命名为 lat 和 lon
	数组形式表示为 map[string]string{lon,lat]
*/
var MappingBanBanUser map[string]interface{} = map[string]interface{}{
	"uid":                     map[string]interface{}{"type": "integer"},
	"app_id":                  map[string]interface{}{"type": "short"},
	"name":                    map[string]interface{}{"type": "text"},
	"sign":                    map[string]interface{}{"type": "text"},
	"icon":                    map[string]interface{}{"type": "text", "index": false},
	"city":                    map[string]interface{}{"type": "text", "index": false},
	"position":                map[string]interface{}{"type": "text", "index": false},
	"birthday":                map[string]interface{}{"type": "integer"},
	"job":                     map[string]interface{}{"type": "integer"},
	"sex":                     map[string]interface{}{"type": "byte"},
	"role":                    map[string]interface{}{"type": "byte"},
	"title":                   map[string]interface{}{"type": "byte"},
	"has_video":               map[string]interface{}{"type": "boolean"},
	"pay_receive_today":       map[string]interface{}{"type": "integer"},
	"interests":               map[string]interface{}{"type": "integer"},
	"god_category":            map[string]interface{}{"type": "integer"},
	"god_num":                 map[string]interface{}{"type": "integer"},
	"god_month_num":           map[string]interface{}{"type": "integer"},
	"god_week_num":            map[string]interface{}{"type": "integer"},
	"god_day_num":             map[string]interface{}{"type": "integer"},
	"god_now_num":             map[string]interface{}{"type": "integer"},
	"god_dateline":            map[string]interface{}{"type": "integer"},
	"god_default_id":          map[string]interface{}{"type": "integer"},
	"god_default_cid":         map[string]interface{}{"type": "integer"},
	"online_status":           map[string]interface{}{"type": "byte"},
	"online_dateline":         map[string]interface{}{"type": "integer"},
	"tag":                     map[string]interface{}{"type": "short"},
	"friend_state":            map[string]interface{}{"type": "short"},
	"online_day":              map[string]interface{}{"type": "integer"},
	"city_code":               map[string]interface{}{"type": "short"},
	"dateline":                map[string]interface{}{"type": "integer"},
	"deleted":                 map[string]interface{}{"type": "byte"},
	"block_un_auther_message": map[string]interface{}{"type": "boolean"},
	"pay_num":                 map[string]interface{}{"type": "integer"},
	"pay_money":               map[string]interface{}{"type": "integer"},
	"pay_room_money":          map[string]interface{}{"type": "integer"},
	"service_score":           map[string]interface{}{"type": "byte"},
	"service_busy":            map[string]interface{}{"type": "integer"},
	"service_pause":           map[string]interface{}{"type": "boolean"},
	"notice_order":            map[string]interface{}{"type": "boolean"},
	"notice_game":             map[string]interface{}{"type": "boolean"},
	"room_rid":                map[string]interface{}{"type": "integer"},
	"room_position":           map[string]interface{}{"type": "byte"},
	"geo":                     map[string]interface{}{"type": "geo_point"},
	"nearby_invisible":        map[string]interface{}{"type": "boolean"},
	"circle_num":              map[string]interface{}{"type": "integer"},
	"photo_num":               map[string]interface{}{"type": "integer"},
	"nearby_day_score":        map[string]interface{}{"type": "double"},
	"recommend_day_score":     map[string]interface{}{"type": "double"},
	"popularity":              map[string]interface{}{"type": "integer"},
	"mentor_level":            map[string]interface{}{"type": "integer"}, //师父等级
	"left_disciple_num":       map[string]interface{}{"type": "integer"}, //剩余可收徒弟的数量
	"gamecard_cids":           map[string]interface{}{"type": "integer"}, //游戏卡id集合
	"exposure": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"exposure_new":     map[string]interface{}{"type": "integer"},
			"exposure_old":     map[string]interface{}{"type": "integer"},
			"follow_new":       map[string]interface{}{"type": "integer"},
			"income_yesterday": map[string]interface{}{"type": "integer"},
			"income_lastweek":  map[string]interface{}{"type": "integer"},
			"income_total":     map[string]interface{}{"type": "integer"},
			"is_peipei":        map[string]interface{}{"type": "boolean"},
			"follows":          map[string]interface{}{"type": "integer"},
			"interests":        map[string]interface{}{"type": "integer"},
			"friends_num":      map[string]interface{}{"type": "integer"},
			"clicked_uids":     map[string]interface{}{"type": "integer"},
		},
	},
}

func UpgradeBanBanUserMapping() {
	//更改 MappingBanBanUser
	//查询所有品类
	cates := []model.XsCategory{}
	_, err := model.Db.Raw("select cid from xs_category where dpath = 2").QueryRows(&cates)
	if err != nil {
		panic(err)
	}

	for _, cate := range cates {
		name := fmt.Sprintf("skill_%d", cate.Cid)
		MappingBanBanUser[name] = map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":          map[string]interface{}{"type": "integer"},
				"cid":         map[string]interface{}{"type": "short"},
				"level":       map[string]interface{}{"type": "short"},
				"disabled":    map[string]interface{}{"type": "boolean"},
				"price":       map[string]interface{}{"type": "integer"},
				"discount":    map[string]interface{}{"type": "byte"},
				"tags":        map[string]interface{}{"type": "integer"},
				"is_default":  map[string]interface{}{"type": "boolean"},
				"num":         map[string]interface{}{"type": "integer"},
				"now_num":     map[string]interface{}{"type": "integer"},
				"day_num":     map[string]interface{}{"type": "integer"},
				"week_num":    map[string]interface{}{"type": "integer"},
				"month_num":   map[string]interface{}{"type": "integer"},
				"cover":       map[string]interface{}{"type": "text", "index": false},
				"audio":       map[string]interface{}{"type": "text", "index": false},
				"sign":        map[string]interface{}{"type": "text", "index": false},
				"description": map[string]interface{}{"type": "text", "index": false},
				"day_credit":  map[string]interface{}{"type": "integer"},
				"week_credit": map[string]interface{}{"type": "integer"},
				"now_credit":  map[string]interface{}{"type": "integer"},
				"closed":      map[string]interface{}{"type": "boolean"},
				"pay_price":   map[string]interface{}{"type": "integer"},
			},
		}
	}
}
