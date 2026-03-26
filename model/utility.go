/*
 * @Descripttion:
 * @Author: Zheng.Jinwei
 * @Date: 2022-08-31 10:43:55
 */
package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/kpango/glg"
)

func GetValueFromParams(values orm.Params) *map[string]string {
	data := map[string]string{}
	for key, value := range values {
		if value == nil {
			data[key] = ""
			continue
		}
		switch val := value.(type) {
		case float64:
			data[key] = strconv.FormatFloat(val, 'f', -1, 64)
		case float32:
			data[key] = strconv.FormatFloat(float64(val), 'f', -1, 64)
		case int:
			data[key] = strconv.Itoa(val)
		case uint:
			data[key] = strconv.Itoa(int(val))
		case int8:
			data[key] = strconv.Itoa(int(val))
		case uint8:
			data[key] = strconv.Itoa(int(val))
		case int16:
			data[key] = strconv.Itoa(int(val))
		case uint16:
			data[key] = strconv.Itoa(int(val))
		case int32:
			data[key] = strconv.Itoa(int(val))
		case uint32:
			data[key] = strconv.Itoa(int(val))
		case int64:
			data[key] = strconv.FormatInt(val, 10)
		case uint64:
			data[key] = strconv.FormatUint(val, 10)
		case string:
			data[key] = val
		case bool:
			if val {
				data[key] = "1"
			} else {
				data[key] = "0"
			}
		case []byte:
			data[key] = string(val)
		default:
			newValue, _ := json.Marshal(value)
			panic(fmt.Errorf("error type, key = %s, val = %s, type = %s", key, string(newValue), reflect.TypeOf(value)))
		}
	}
	return &data
}

func getIntFromParams(values orm.Params, name string) int {
	for key, value := range values {
		if value == nil {
			continue
		}

		if key != name {
			continue
		}

		uid, _ := strconv.Atoi(fmt.Sprintf("%v", value))
		return uid
	}

	return 0
}

func BatchGetUnionChatRoomUidAll(ridIndex, limit int) (map[int]int, int, error) {
	res := []orm.Params{}
	rawSql := fmt.Sprintf("select * from xs_chatroom where property='union' and rid > %d limit %d", ridIndex, limit)
	_, err := MasterDb.Raw(rawSql).Values(&res)

	if err != nil {
		glg.Error("BatchGetUnionChatRoomUidAll err:", err)
		return nil, ridIndex, err
	}

	retMap := make(map[int]int, limit)
	for _, value := range res {
		tmpRid := getIntFromParams(value, "rid")
		uid := getIntFromParams(value, "uid")
		if tmpRid > 0 && uid > 0 {
			retMap[tmpRid] = uid / 100
		}
		if ridIndex < tmpRid {
			ridIndex = tmpRid
		}
	}

	glg.Debug("BatchGetUnionChatRoomUidAll retMap:", retMap)
	return retMap, ridIndex, nil
}

func BatchGetChatRoomUidAll(ridIndex, limit int) (map[int]int, int, error) {
	res := []orm.Params{}
	rawSql := fmt.Sprintf("select * from xs_chatroom where rid > %d order by rid limit %d", ridIndex, limit)
	_, err := MasterDb.Raw(rawSql).Values(&res)

	if err != nil {
		glg.Error("BatchGetChatRoomUidAll err:", err)
		return nil, ridIndex, err
	}

	retMap := make(map[int]int, limit)
	for _, value := range res {
		tmpRid := getIntFromParams(value, "rid")
		uid := getIntFromParams(value, "uid")
		if tmpRid > 0 && uid > 0 {
			retMap[tmpRid] = uid
		}
		if ridIndex < tmpRid {
			ridIndex = tmpRid
		}
	}

	glg.Debug("BatchGetChatRoomUidAll retMap:", retMap)
	return retMap, ridIndex, nil
}

func batchGetBigareaIdByUid(uidList []int) (map[int]int, error) {
	if len(uidList) == 0 {
		return nil, nil
	}
	uidStr := strings.Replace(strings.Trim(fmt.Sprint(uidList), "[]"), " ", ",", -1)

	rawSql := fmt.Sprintf("select * from xs_user_bigarea where uid in (%s)", uidStr)
	res := []orm.Params{}
	_, err := MasterDb.Raw(rawSql).Values(&res)

	if err != nil {
		glg.Error("batchGetBigareaIdByUid err:", err)
		return nil, err
	}

	//glg.Debug("batchGetBigareaIdByUid res:", res)
	retMap := make(map[int]int, len(uidList))
	for _, value := range res {
		rid := getIntFromParams(value, "uid")
		uid := getIntFromParams(value, "bigarea_id")
		if rid > 0 && uid > 0 {
			retMap[rid] = uid
		}
	}

	glg.Debug("batchGetBigareaIdByUid retMap:", retMap)
	return retMap, nil
}

func batchGetUidByRid(ridList []int) (map[int]int, error) {
	if len(ridList) == 0 {
		return nil, nil
	}
	ridStr := strings.Replace(strings.Trim(fmt.Sprint(ridList), "[]"), " ", ",", -1)

	res := []orm.Params{}
	rawSql := fmt.Sprintf("select * from xs_chatroom where rid in (%s)", ridStr)
	_, err := MasterDb.Raw(rawSql).Values(&res)

	if err != nil {
		glg.Error("BatchGetUidByRid err:", err)
		return nil, err
	}

	retMap := make(map[int]int, len(ridList))
	for _, value := range res {
		rid := getIntFromParams(value, "rid")
		uid := getIntFromParams(value, "uid")
		if rid > 0 && uid > 0 {
			retMap[rid] = uid / 100
		}
	}

	glg.Debug("BatchGetUidByRid retMap:", retMap)

	return retMap, nil
}

func BatchGetRidBigAreaId(ridList []int) (map[int]int, error) {
	ridUidMap, err := batchGetUidByRid(ridList)
	if err != nil {
		return nil, err
	}

	_, uidList := getIntMapList(ridUidMap)
	bigAreaIdMap, err := BatchGetBigAreaIdByUids(uidList)
	if err != nil {
		return nil, err
	}

	ridAreaIdmap := make(map[int]int, len(ridUidMap))
	for rid, uid := range ridUidMap {
		bigAreaId := bigAreaIdMap[uid]
		if bigAreaId > 0 {
			ridAreaIdmap[rid] = bigAreaId
		}
	}

	return ridAreaIdmap, nil
}

func BatchGetRidBigAreaIdAll(ridIndex, limit int) (map[int]int, int, int, error) {
	ridIdMap, maxRid, err := BatchGetChatRoomUidAll(ridIndex, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	_, uidList := getIntMapList(ridIdMap)
	bigAreaIdMap, err := BatchGetBigAreaIdByUids(uidList)
	if err != nil {
		return nil, maxRid, len(ridIdMap), err
	}

	ridAreaIdmap := make(map[int]int, len(ridIdMap))
	for rid, uid := range ridIdMap {
		bigAreaId := bigAreaIdMap[uid]
		if bigAreaId > 0 {
			ridAreaIdmap[rid] = bigAreaId
		}
	}

	return ridAreaIdmap, maxRid, len(ridIdMap), nil
}

func BatchGetBigAreaIdByUids(uidList []int) (map[int]int, error) {
	bigAreaIdMap, err := batchGetBigareaIdByUid(uidList)
	if err != nil {
		return nil, err
	}

	return bigAreaIdMap, nil
}

func getIntMapList(srcMap map[int]int) ([]int, []int) {
	ridList := make([]int, 0, len(srcMap))

	hasUidMap := make(map[int]bool)
	for rid, uid := range srcMap {
		hasUidMap[uid] = true
		ridList = append(ridList, rid)
	}

	uidList := make([]int, 0, len(hasUidMap))
	for uid := range hasUidMap {
		uidList = append(uidList, uid)
	}

	return ridList, uidList
}

func GetBigAreaIdByUid(uid int) (bigareaId int, err error) {
	rawSql := fmt.Sprintf("select * from xs_user_bigarea where uid = %v", uid)
	res := []orm.Params{}
	_, err = MasterDb.Raw(rawSql).Values(&res)

	if err != nil {
		glg.Error("GetBigAreaIdByUid err:", err)
		return
	}

	for _, value := range res {
		bigareaId = getIntFromParams(value, "bigarea_id")
	}

	return
}
