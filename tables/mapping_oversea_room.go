/*
 * @Descripttion:
 * @Author: Zheng.Jinwei
 * @Date: 2022-08-31 10:43:55
 */
package tables

var MappingOverseaRoom = map[string]interface{}{
	"rid":                 map[string]interface{}{"type": "integer"},
	"uid":                 map[string]interface{}{"type": "integer"},
	"app_id":              map[string]interface{}{"type": "short"},
	"name":                map[string]interface{}{"type": "text"},
	"icon":                map[string]interface{}{"type": "text", "index": false},
	"room_icon":           map[string]interface{}{"type": "text", "index": false},
	"password":            map[string]interface{}{"type": "boolean"},
	"prefix":              map[string]interface{}{"type": "text", "index": false},
	"property":            map[string]interface{}{"type": "keyword"},
	"type":                map[string]interface{}{"type": "keyword"},
	"types":               map[string]interface{}{"type": "keyword"},
	"dateline":            map[string]interface{}{"type": "integer"},
	"deleted":             map[string]interface{}{"type": "byte"},
	"online_num":          map[string]interface{}{"type": "integer"},
	"emperor_time":        map[string]interface{}{"type": "integer"},
	"reception_uid":       map[string]interface{}{"type": "integer"},
	"fans_num":            map[string]interface{}{"type": "integer"},
	"boy_num":             map[string]interface{}{"type": "integer"},
	"girl_num":            map[string]interface{}{"type": "integer"},
	"boss_uid":            map[string]interface{}{"type": "integer"},
	"uname":               map[string]interface{}{"type": "text", "index": false},
	"utitle":              map[string]interface{}{"type": "integer"},
	"lock_end":            map[string]interface{}{"type": "integer"},
	"mic_num":             map[string]interface{}{"type": "integer"},
	"show_start":          map[string]interface{}{"type": "integer"},
	"weight":              map[string]interface{}{"type": "integer"},
	"room_sex":            map[string]interface{}{"type": "byte"},
	"factory_type":        map[string]interface{}{"type": "keyword"},
	"sex":                 map[string]interface{}{"type": "byte"},
	"state":               map[string]interface{}{"type": "integer"}, //0：正常，1：隐藏，2：仅好友可见，3：仅粉丝可见
	"paier":               map[string]interface{}{"type": "integer"}, //是否显示老板位
	"show":                map[string]interface{}{"type": "keyword"}, //房间外显标签 -- 已不用
	"settlement_channel":  map[string]interface{}{"type": "keyword"}, //结算频道
	"pin":                 map[string]interface{}{"type": "integer"}, //房间置顶
	"fixed_tag_id":        map[string]interface{}{"type": "text"},    //房间外显标签 -- 已玩坏
	"fixed_tag_id_new":    map[string]interface{}{"type": "integer"}, //房间外显标签
	//
	"room_available_seat": map[string]interface{}{"type": "integer"},
	"room_occupied_seat":  map[string]interface{}{"type": "integer"},
	"room_total_seat":     map[string]interface{}{"type": "integer"},
	//
	"call_state":          map[string]interface{}{"type": "integer"},
	"call_time":          map[string]interface{}{"type": "integer"},
	//
	"min_age":          map[string]interface{}{"type": "float"},
	"max_age":          map[string]interface{}{"type": "float"},
	"avg_age":          map[string]interface{}{"type": "float"},
	//
	"blocked": map[string]interface{}{"type": "boolean"},
	//
	"room_gaming_level": map[string]interface{}{"type": "integer"},
	//
	"recruit_uid": map[string]interface{}{"type": "integer"},
	"top_category": map[string]interface{}{"type": "integer"},
	"recruit_status": map[string]interface{}{"type": "integer"},
	"recruit_version": map[string]interface{}{"type": "integer"},
	//
	"gaming_name":  map[string]interface{}{"type": "keyword"},
	"min_players":  map[string]interface{}{"type": "integer"},
	"gaming_state": map[string]interface{}{"type": "integer"},

	// 狼人杀等级
	"wolf_level": map[string]interface{}{"type": "integer"},

	//增加匹配相关的
	"game_state":        map[string]interface{}{"type": "keyword"}, //游戏状态，使用字符串，通用
	"game_state_v2":     map[string]interface{}{"type": "integer"}, //游戏状态，数字，直接来自游戏状态表
	"game":              map[string]interface{}{"type": "keyword"}, //游戏类型
	"game_state_create": map[string]interface{}{"type": "integer"}, //游戏状态记录创建时间
	"game_state_update": map[string]interface{}{"type": "integer"}, //游戏状态记录最后更新时间
	"game_version":      map[string]interface{}{"type": "integer"}, //游戏的房间版本号

	//海外特有的
	"area":       map[string]interface{}{"type": "keyword"},
	"language":   map[string]interface{}{"type": "keyword"},
	"real":       map[string]interface{}{"type": "integer"},
	"bigarea_id": map[string]interface{}{"type": "byte"},

	"base": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"new_male_stay_seconds":   map[string]interface{}{"type": "integer"},
			"new_female_stay_seconds": map[string]interface{}{"type": "integer"},
			"old_male_stay_seconds":   map[string]interface{}{"type": "integer"},
			"old_female_stay_seconds": map[string]interface{}{"type": "integer"},
			"new_male_pay_percent":    map[string]interface{}{"type": "float"},
			"new_female_pay_percent":  map[string]interface{}{"type": "float"},
			"old_male_pay_percent":    map[string]interface{}{"type": "float"},
			"old_female_pay_percent":  map[string]interface{}{"type": "float"},
		},
	},
	"round_real":       map[string]interface{}{"type": "integer"},

	"tags":                 map[string]interface{}{"type": "keyword"},
	"pk_state":             map[string]interface{}{"type": "long"},
	"team_pk_state":        map[string]interface{}{"type": "long"},
	"red_packet_post_time": map[string]interface{}{"type": "long"},
	"link_mic_status":      map[string]interface{}{"type": "long"},
	"online_robot_num":     map[string]interface{}{"type": "long"},

	"follow_num":          map[string]interface{}{"type": "long"},
	"follow_num_dateline": map[string]interface{}{"type": "long"},

	"link_mic_num":          map[string]interface{}{"type": "long"},
	"link_mic_num_dateline": map[string]interface{}{"type": "long"},

	"comment_num":          map[string]interface{}{"type": "long"},
	"comment_num_dateline": map[string]interface{}{"type": "long"},

	"money":          map[string]interface{}{"type": "long"},
	"money_dateline": map[string]interface{}{"type": "long"},

	"family_id":                map[string]interface{}{"type": "long"},
	"boom_rocket_lv":           map[string]interface{}{"type": "long"},
	"boom_rocket_failure_time": map[string]interface{}{"type": "long"},
}
