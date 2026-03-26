package main

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/conf"
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"github.com/olachat/banban_server/cli_mysql_pipeline/tables"
	"github.com/kpango/glg"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	//解析命令行参数
	var names string
	var env string
	flag.StringVar(&env, "env", "", "操作哪个配置文件数据")
	flag.StringVar(&names, "names", "", "更新哪个索引的mapping")
	flag.Parse()

	if len(names) == 0 {
		panic(fmt.Errorf("error index name"))
	}

	conf.ParseConfig(env)

	model.InitDb()

	tables.UpgradeBanBanUserMapping()
	tables.UpgradeOverseaUserMapping()

	//聊天室的 room tables.MappingBanBanRoom
	//用户的 user tables.MappingBanBanUser
	//大神的 user_god tables.MappingBanBanUser

	//海外的 room_new user_new user_god_new

	array := strings.Split(names, ",")
	for _, name := range array {
		switch name {
		case "room":
			updateMapping(name, tables.MappingBanBanRoom)
		case "user_god_new":
			updateMapping(name, tables.MappingBanBanUser)
		case "user_new":
			updateMapping(name, tables.MappingBanBanUser)
		case "oversea.room_new":
			updateMapping("room_new", tables.MappingOverseaRoom)
		case "oversea.user_new":
			updateMapping("user_new", tables.MappingOverseaUser)
		case "oversea.user_god_new":
			updateMapping("user_god_new", tables.MappingOverseaUser)
		default:
			panic(fmt.Errorf("error name %s", name))
		}
	}
}

func updateMapping(name string, mapping map[string]interface{}) {
	glg.Info("update mapping", name)
	post := map[string]interface{}{
		"properties": mapping,
	}
	curl(
		"PUT",
		fmt.Sprintf("%s/_mapping/default", name),
		post,
	)
}

func curl(method string, path string, jsonData map[string]interface{}) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", conf.EsHost, path), bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if len(conf.EsName) > 0 && len(conf.EsPass) > 0 {
		req.SetBasicAuth(conf.EsName, conf.EsPass)
	}

	resp, err := client.Do(req)

	if jsonData != nil {
		glg.Info("update mapping json \n", string(body))
	}

	output, _ := ioutil.ReadAll(resp.Body)
	glg.Info(string(output))
	if err != nil {
		panic(err)
	}
	return nil
}
