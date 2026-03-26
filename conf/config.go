package conf

import (
	"fmt"
	"io/ioutil"
	"path"

	utility "github.com/olachat/banban_server/common"

	yaml "gopkg.in/yaml.v2"
)

var MysqlDebug bool = false

var MysqlAddr string
var MysqlMasterAddr string

var MqAddr string
var MqTopic string
var MqGroup string
var MqUserName string
var MqPassword string

var NsqAddr string
var NsqTopic string
var NsqChannel string

var EsHost string
var EsName string
var EsPass string

var IsDev bool = true
var Discovery string

// NsqPublisherConf NSQ Sink 配置（nil 表示不启用 NSQ Sink）
type NsqPublisherConf struct {
	Pools     map[string][]string `yaml:"pools"`      // poolName → NSQ 地址列表
	Topics    map[string]string   `yaml:"topics"`     // topicName → poolName
	Tables    map[string][]string `yaml:"tables"`     // tableName → topic 列表
	NewTables []string            `yaml:"new_tables"` // 使用 JSON 格式的表（其余用 PHP 序列化）
	Debug     bool                `yaml:"debug"`
}

type Conf struct {
	MysqlDebug      bool              `yaml:"mysql_debug"`
	MysqlAddr       string            `yaml:"mysql_addr"`
	MysqlMasterAddr string            `yaml:"mysql_master_addr"`
	EsHost          string            `yaml:"es_host"`
	EsName          string            `yaml:"es_name"`
	EsPass          string            `yaml:"es_pass"`
	MqAddr          string            `yaml:"mq_addr"`
	MqTopic         string            `yaml:"mq_topic"`
	MqGroup         string            `yaml:"mq_group"`
	MqUserName      string            `yaml:"mq_user_name"`
	MqPassword      string            `yaml:"mq_password"`
	NsqAddr         string            `yaml:"nsq_addr"`
	NsqTopic        string            `yaml:"nsq_topic"`
	NsqChannel      string            `yaml:"nsq_channel"`
	ENV             string            `yaml:"env"`
	Discovery       string            `yaml:"discovery"`
	NsqPublisher    *NsqPublisherConf `yaml:"nsq_publisher,omitempty"`
}

// ParseConfig 读取并解析配置文件，设置全局变量，返回完整 Conf（供调用方访问 NsqPublisher 等扩展字段）
func ParseConfig(name string) *Conf {
	fileName := path.Join(
		utility.ExecPath(),
		fmt.Sprintf("%s.yaml", name),
	)
	fmt.Println("config fileName:", fileName)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	conf := &Conf{}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		panic(err)
	}

	MysqlDebug = conf.MysqlDebug
	MysqlAddr = conf.MysqlAddr
	MysqlMasterAddr = conf.MysqlMasterAddr

	EsHost = conf.EsHost
	EsName = conf.EsName
	EsPass = conf.EsPass

	MqAddr = conf.MqAddr
	MqTopic = conf.MqTopic
	MqGroup = conf.MqGroup
	MqUserName = conf.MqUserName
	MqPassword = conf.MqPassword

	NsqAddr = conf.NsqAddr
	NsqTopic = conf.NsqTopic
	NsqChannel = conf.NsqChannel

	IsDev = conf.ENV == "dev"
	Discovery = conf.Discovery

	fmt.Printf("%+v\n", conf)
	return conf
}
