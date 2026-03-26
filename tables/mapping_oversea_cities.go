package tables

import (
	"encoding/json"
)

var MappingCities map[string]interface{} = map[string]interface{}{
	"state_id":   map[string]interface{}{"type": "integer", "index": false},
	"country_id": map[string]interface{}{"type": "integer", "index": false},
	"location":   map[string]interface{}{"type": "geo_point"},
	"ar":         map[string]interface{}{"type": "text"},
	"en":         map[string]interface{}{"type": "text"},
	"id":         map[string]interface{}{"type": "text"},
	"ko":         map[string]interface{}{"type": "text"},
	"ms":         map[string]interface{}{"type": "text"},
	"th":         map[string]interface{}{"type": "text"},
	"tr":         map[string]interface{}{"type": "text"},
	"vi":         map[string]interface{}{"type": "text"},
	"zh_cn":      map[string]interface{}{"type": "text"},
	"zh_tw":      map[string]interface{}{"type": "text"},
}

type translation struct {
	Ar    string `json:"ar"`
	En    string `json:"en"`
	Id    string `json:"id"`
	Ko    string `json:"ko"`
	Ms    string `json:"ms"`
	Th    string `json:"th"`
	Tr    string `json:"tr"`
	Vi    string `json:"vi"`
	Zh_cn string `json:"zh_cn"`
	Zh_tw string `json:"zh_tw"`
}

func CityNameFactoryGivenLang(lang string) func(map[string]string) interface{} {
	return func(origin map[string]string) interface{} {
		var (
			res                 string
			cityName, stateName translation
		)

		if err := json.Unmarshal([]byte(origin["translation"]), &cityName); err == nil {
			switch lang {
			case "ar":
				res = cityName.Ar
			case "en":
				res = cityName.En
			case "id":
				res = cityName.Id
			case "ko":
				res = cityName.Ko
			case "ms":
				res = cityName.Ms
			case "th":
				res = cityName.Th
			case "tr":
				res = cityName.Tr
			case "vi":
				res = cityName.Vi
			case "zh_cn":
				res = cityName.Zh_cn
			case "zh_tw":
				res = cityName.Zh_tw
			}
		}
		if err := json.Unmarshal([]byte(origin["state_translation"]), &stateName); err == nil {
			switch lang {
			case "ar":
				res += " " + stateName.Ar
			case "en":
				res += " " + stateName.En
			case "id":
				res += " " + stateName.Id
			case "ko":
				res += " " + stateName.Ko
			case "ms":
				res += " " + stateName.Ms
			case "th":
				res += " " + stateName.Th
			case "tr":
				res += " " + stateName.Tr
			case "vi":
				res += " " + stateName.Vi
			case "zh_cn":
				res += " " + stateName.Zh_cn
			case "zh_tw":
				res += " " + stateName.Zh_tw
			}
		}
		return res
	}
}
