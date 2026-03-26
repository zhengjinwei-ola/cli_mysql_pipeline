package tables

var MappingBanBanRoom map[string]interface{} = map[string]interface{}{
	"rid":            map[string]interface{}{"type": "integer"},
	"uid":            map[string]interface{}{"type": "integer"},
	"app_id":         map[string]interface{}{"type": "short"},
	"name":           map[string]interface{}{"type": "text"},
	"icon":           map[string]interface{}{"type": "text", "index": false},
	"room_icon":      map[string]interface{}{"type": "text", "index": false},
	"password":       map[string]interface{}{"type": "boolean"},
	"prefix":         map[string]interface{}{"type": "text", "index": false},
	"property":       map[string]interface{}{"type": "keyword"},
	"type":           map[string]interface{}{"type": "keyword"},
	"types":          map[string]interface{}{"type": "keyword"},
	"dateline":       map[string]interface{}{"type": "integer"},
	"deleted":        map[string]interface{}{"type": "byte"},
	"online_num":     map[string]interface{}{"type": "integer"},
	"emperor_time":   map[string]interface{}{"type": "integer"},
	"reception_uid":  map[string]interface{}{"type": "integer"},
	"fans_num":       map[string]interface{}{"type": "integer"},
	"boy_num":        map[string]interface{}{"type": "integer"},
	"girl_num":       map[string]interface{}{"type": "integer"},
	"boss_uid":       map[string]interface{}{"type": "integer"},
	"boss_icon":      map[string]interface{}{"type": "text"},
	"uname":          map[string]interface{}{"type": "text", "index": false},
	"utitle":         map[string]interface{}{"type": "integer"},
	"live_state":     map[string]interface{}{"type": "text"}, //直播游戏房房间直播状态。默认为空(pending, playing, ending)
	"factory_type":   map[string]interface{}{"type": "keyword"},
	"fixed_tag_id":   map[string]interface{}{"type": "integer"},
	"module_id":      map[string]interface{}{"type": "integer"},
	"settle_channel": map[string]interface{}{"type": "keyword"},

	//增加匹配相关的
	"game_state":        map[string]interface{}{"type": "keyword"}, //游戏状态，使用字符串，通用
	"game_state_v2":     map[string]interface{}{"type": "integer"}, //游戏状态，数字，直接来自游戏状态表
	"game":              map[string]interface{}{"type": "keyword"}, //游戏类型
	"game_state_create": map[string]interface{}{"type": "integer"}, //游戏状态记录创建时间
	"game_state_update": map[string]interface{}{"type": "integer"}, //游戏状态记录最后更新时间
	"game_version":      map[string]interface{}{"type": "integer"}, //游戏的房间版本号

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
}
