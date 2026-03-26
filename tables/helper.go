package tables

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
)

type OpRowTable struct {
}

func (o OpRowTable) InitHook(docId int64, data *map[string]interface{}) {

}

func (o OpRowTable) InitHookBefore(data *map[string]string) {

}

func (o OpRowTable) Format(
	origin map[string]string,
	before map[string]string,
	docField string,
	fields *[]model.EsField,
	op model.OpType,
) (int64, map[string]interface{}, model.OpType) {
	value, ok := origin[docField]
	if !ok {
		return 0, nil, op
	}
	docId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, nil, op
	}
	if op == model.OpType_Update {
		afterImages, errorAfter := o.format(origin, fields)
		beforeImages, errorBefore := o.format(before, fields)
		// fmt.Println("[debug-after]", afterImages)
		// fmt.Println("[debug-before]", beforeImages)
		if errorBefore != nil || errorAfter != nil {
			fmt.Println("errorBefore", errorBefore)
			fmt.Println("errorAfter", errorAfter)
			return 0, nil, op
		}
		fmt.Println("[debug-after]", afterImages)
		fmt.Println("[debug-before]", beforeImages)
		data := map[string]interface{}{}
		for key, val := range afterImages {
			// 修复es没数据问题 打开下面两行 更新数据
			//data[key] = val
			//continue
			v, ok := beforeImages[key]
			if !ok {
				data[key] = val
				continue
			}
			if !reflect.DeepEqual(val, v) {
				data[key] = val
			}
		}
		if len(data) == 0 {
			return 0, nil, op
		}
		return docId, data, op
	} else {
		value, err := o.format(origin, fields)
		if err != nil {
			fmt.Println("error", err)
			return 0, nil, op
		}
		return docId, value, op
	}
}

func (o OpRowTable) format(
	origin map[string]string,
	fields *[]model.EsField,
) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for _, field := range *fields {
		target := field.Name
		if len(field.Target) > 0 {
			target = field.Target
		}
		if field.Type == model.EsType_Func {
			data[target] = field.Func(origin)
			continue
		}
		str, ok := origin[field.Name]
		if !ok {
			return nil, fmt.Errorf("error field %s in data %v", field.Name, origin)
		}
		switch field.Type {
		case model.EsType_Number:
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				v = 0
			}
			data[target] = v
		case model.EsType_Float:
			v, err := strconv.ParseFloat(str, 64)
			if err != nil {
				v = 0
			}
			data[target] = v
		case model.EsType_Text:
			data[target] = str
		case model.EsType_TureOrFalse:
			if len(str) > 0 {
				switch str {
				case model.True_Upper, model.True_Lower:
					data[target] = true
				case model.False_Upper, model.False_Lower:
					data[target] = false
				default:
					v, err := strconv.Atoi(str)
					if err != nil {
						v = 0
					}
					data[target] = v > 0
				}
			} else {
				data[target] = false
			}
		case model.EsType_SetNumber:
			if len(str) == 0 {
				data[target] = []int64{}
			} else {
				v := []int64{}
				array := strings.Split(str, ",")
				for _, s := range array {
					a, err := strconv.ParseInt(s, 10, 64)
					if err == nil {
						v = append(v, a)
					}
				}
				data[target] = v
			}
		case model.EsType_SetString:
			if len(str) == 0 {
				data[target] = []string{}
			} else {
				data[target] = strings.Split(str, ",")
			}
		}
	}
	return data, nil
}
