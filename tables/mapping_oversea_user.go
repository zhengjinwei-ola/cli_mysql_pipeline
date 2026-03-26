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
var MappingOverseaUser = map[string]interface{}{
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
	"popularity":              map[string]interface{}{"type": "integer"},
	"mentor_level":            map[string]interface{}{"type": "integer"}, //师父等级
	"left_disciple_num":       map[string]interface{}{"type": "integer"}, //剩余可收徒弟的数量
	"meet_setting_sex":        map[string]interface{}{"type": "integer"}, //遇见-偏好-性别
	"personal_tag":            map[string]interface{}{"type": "integer"}, // 个人标签
	"interest_tag":            map[string]interface{}{"type": "integer"}, // 兴趣标签
	"bigarea_id":              map[string]interface{}{"type": "byte"},    // 大区
	//漏掉的
	"room_property": map[string]interface{}{"type": "keyword"},
	"room_type":     map[string]interface{}{"type": "keyword"},
	"room_types":    map[string]interface{}{"type": "keyword"},
	"country":       map[string]interface{}{"type": "keyword"},
	"country_code":  map[string]interface{}{"type": "keyword"},
	"language":      map[string]interface{}{"type": "keyword"},
	"room_game":     map[string]interface{}{"type": "keyword"},
	"room_weight":   map[string]interface{}{"type": "integer"},
	"room_name":     map[string]interface{}{"type": "text", "index": false},

	// 用户位置相关
	"viewing_page":   map[string]interface{}{"type": "integer"},
	// 游戏自动推送匹配优化相关
	"mute_sys_game_prompts":   map[string]interface{}{"type": "boolean"},
	"mute_stranger_game_invite":   map[string]interface{}{"type": "boolean"},
	"mute_sys_game_prompts_till":   map[string]interface{}{"type": "integer"},
	"mute_stranger_game_invite_till":   map[string]interface{}{"type": "integer"},
	"sys_game_match_prompted_to":   map[string]interface{}{"type": "boolean"},
	"sys_game_match_expire_at":   map[string]interface{}{"type": "integer"},
	// guess
	"recent_guess_total":   map[string]interface{}{"type": "integer"},
	"recent_guess_score":   map[string]interface{}{"type": "integer"},
	"recent_guess_thumbs_up":   map[string]interface{}{"type": "integer"},
	"recent_guess_last_played":   map[string]interface{}{"type": "integer"},
	// under
	"recent_under_total":   map[string]interface{}{"type": "integer"},
	"recent_under_populace":   map[string]interface{}{"type": "integer"},
	"recent_under_under":   map[string]interface{}{"type": "integer"},
	"recent_under_last_played":   map[string]interface{}{"type": "integer"},
	// wolf
	"recent_wolf_total":   map[string]interface{}{"type": "integer"},
	"recent_wolf_win":   map[string]interface{}{"type": "integer"},
	"recent_wolf_goodWin":   map[string]interface{}{"type": "integer"},
	"recent_wolf_level":   map[string]interface{}{"type": "integer"},
	"recent_wolf_exp":   map[string]interface{}{"type": "integer"},
	"recent_wolf_last_played":   map[string]interface{}{"type": "integer"},
	// unity - ludo
	"recent_ludo_round":   map[string]interface{}{"type": "integer"},
	"recent_ludo_champion":   map[string]interface{}{"type": "integer"},
	"recent_ludo_score":   map[string]interface{}{"type": "integer"},
	"recent_ludo_last_played":   map[string]interface{}{"type": "integer"},
	"recent_ludo_dateline":   map[string]interface{}{"type": "integer"},
	// unity - carrom
	"recent_carrom_round":   map[string]interface{}{"type": "integer"},
	"recent_carrom_champion":   map[string]interface{}{"type": "integer"},
	"recent_carrom_score":   map[string]interface{}{"type": "integer"},
	"recent_carrom_last_played":   map[string]interface{}{"type": "integer"},
	"recent_carrom_dateline":   map[string]interface{}{"type": "integer"},
	// unity -billiards
	"recent_billiards_round":   map[string]interface{}{"type": "integer"},
	"recent_billiards_champion":   map[string]interface{}{"type": "integer"},
	"recent_billiards_score":   map[string]interface{}{"type": "integer"},
	"recent_billiards_last_played":   map[string]interface{}{"type": "integer"},
	"recent_billiards_dateline":   map[string]interface{}{"type": "integer"},
	// puzzle
	"recent_puzzle_total":   map[string]interface{}{"type": "integer"},
	"recent_puzzle_score":   map[string]interface{}{"type": "integer"},
	"recent_puzzle_last_played":   map[string]interface{}{"type": "integer"},

	// 曝光度
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

func UpgradeOverseaUserMapping() {
	//更改 MappingBanBanUser
	//查询所有品类
	cates := []model.XsCategory{}
	_, err := model.Db.Raw("select cid from xs_category where dpath = 2").QueryRows(&cates)
	if err != nil {
		panic(err)
	}

	for _, cate := range cates {
		if cate.Cid == 372 || cate.Cid == 374 || cate.Cid == 395 || cate.Cid == 397 {
			continue
		}
		name := fmt.Sprintf("skill_%d", cate.Cid)
		MappingOverseaUser[name] = map[string]interface{}{
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
