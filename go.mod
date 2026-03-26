module github.com/olachat/banban_server/cli_mysql_pipeline

go 1.16

require (
	github.com/Shopify/sarama v1.29.1
	github.com/actgardner/gogen-avro/v7 v7.3.1
	github.com/astaxie/beego v1.12.3
	github.com/dolthub/go-mysql-server v0.10.0
	github.com/elliotchance/phpserialize v1.3.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gogf/gf v1.16.9
	github.com/golang/mock v1.6.0
	github.com/google/wire v0.5.0
	github.com/kpango/glg v1.6.4
	github.com/nsqio/go-nsq v1.0.8
	github.com/nsqio/nsq v1.2.1
	github.com/olachat/banban_server v0.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.4
	github.com/syyongx/php2go v0.9.4
	github.com/tidwall/gjson v1.18.0
	github.com/yvasiyarov/php_session_decoder v0.0.0-20180803065642-a065a3b0b7d1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/olachat/banban_server => ../psl-be-banban_server
	github.com/smallnest/rpcx => github.com/olaola-chat/rpcx v0.1.1
)
