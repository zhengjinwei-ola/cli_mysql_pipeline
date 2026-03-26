package sink

import "github.com/olachat/banban_server/cli_mysql_pipeline/model"

// Sink 是消息 Sink 的统一接口，App Router 通过此接口向各 Sink 分发消息
type Sink interface {
	// Name 返回 Sink 标识，用于日志
	Name() string

	// Tables 返回该 Sink 关心的 table 列表（"db.table" 格式）
	// App Router 用此列表构建路由表
	Tables() []string

	// Send 发送一条消息到该 Sink（写入内部 channel，非阻塞）
	Send(row *model.OriginRow)

	// Start 启动 Sink 内部 goroutine
	Start()

	// Stop 优雅停止
	Stop()
}
