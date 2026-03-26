package song_index

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/index"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
)

func GetIndex() *index.Index {
	return &index.Index{
		Name:             "song",
		Receiver:         make(chan *model.OriginRow, 1),
		Mapping:          tables.MappingOverseaSong,
		NumberOfShards:   1,
		NumberOfReplicas: 2,
		Ops: []tables.OpRowBase{
			{
				Db:       "xianshi",
				Name:     "xs_ktv_song",
				DocField: "id",
				PkField:  "id",
				Flow:     model.FlowType_Main,
				Fields: &[]model.EsField{
					{Name: "id", Type: model.EsType_Number},
					{Name: "name", Type: model.EsType_Text},
					{Name: "photo", Type: model.EsType_Text},
					{Name: "singer_id", Type: model.EsType_Number},
					{Name: "singer_name", Type: model.EsType_Text},
					{Name: "original_mp3", Type: model.EsType_Text},
					{Name: "size", Type: model.EsType_Number},
					{Name: "playtime", Type: model.EsType_Number},
					{Name: "hq_music", Type: model.EsType_Text},
					{Name: "hq_size", Type: model.EsType_Number},
					{Name: "hq_playtime", Type: model.EsType_Number},
					{Name: "brc", Type: model.EsType_Text},
					{Name: "uploader_uid", Type: model.EsType_Number},
					{Name: "uploader_name", Type: model.EsType_Text},
					{Name: "uploader_photo", Type: model.EsType_Text},
					{Name: "tag", Type: model.EsType_SetString},
					{Name: "type", Type: model.EsType_SetNumber},
					{Name: "status", Type: model.EsType_Number},
					{Name: "language", Type: model.EsType_Text},
					{Name: "hq_status", Type: model.EsType_Number},
					{Name: "dateline", Type: model.EsType_Number},
					{Name: "updateline", Type: model.EsType_Number},
				},
				Wrap: tables.OpRowTable{},
			},
		},
	}
}
