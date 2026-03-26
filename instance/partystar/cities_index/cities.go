package cities_index

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"strconv"
)

func GetIndex() *index.Index {
	return &index.Index{
		Name:             "cities_index",
		Receiver:         make(chan *model.OriginRow, 1),
		Mapping:          tables.MappingCities,
		NumberOfShards:   1,
		NumberOfReplicas: 1,
		Ops: []tables.OpRowBase{
			{
				Db:       "xianshi",
				Name:     "xs_cities",
				DocField: "id",
				PkField:  "id",
				SqlTemplate: `select s.translation as state_translation, s.id as state_id, co.id as country_id, c.*
				from xs_states as s inner join xs_cities as c on s.id = c.state_id inner join xs_country as co
				on co.id = s.country_id where c.id >= ? and c.id <= ?`,
				Flow: model.FlowType_Main,
				Fields: &[]model.EsField{
					{Name: "state_id", Type: model.EsType_Number},
					{Name: "country_id", Type: model.EsType_Number},
					{
						Name: "location",
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
					{
						Name: "ar",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("ar"),
					},
					{
						Name: "en",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("en"),
					},
					{
						Name: "id",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("id"),
					},
					{
						Name: "ko",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("ko"),
					},
					{
						Name: "ms",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("ms"),
					},
					{
						Name: "th",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("th"),
					},
					{
						Name: "tr",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("tr"),
					},
					{
						Name: "vi",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("vi"),
					},
					{
						Name: "zh_cn",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("zh_cn"),
					},
					{
						Name: "zh_tw",
						Type: model.EsType_Func,
						Func: tables.CityNameFactoryGivenLang("zh_tw"),
					},
				},
				Wrap: tables.OpRowTable{},
			},
		},
	}
}
