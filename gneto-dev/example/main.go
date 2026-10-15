package main

import (
	"gneto"
	"gneto/internal/logger"
	"gneto/pkg/config"
	"gneto/pkg/processor"
	"gneto/pkg/router"
	"gneto/types"
)

// 模拟 Protobuf 消息
type LoginReq struct{ Code int }
type LoginRes struct{ Code int }

const (
	MsgId_Login    = 1001
	MsgId_LoginRes = 1002
)

func main() {
	// 示例代码与基准测试 (Example code and benchmarks)

	r := router.NewRouter()

	// 初始化服务器
	p := processor.NewProcessor(config.WithPayloadType(types.JSON))
	server := gneto.NewServer()
	err := server.Init(
		r.RouterFunc(),
		p,
		config.WithPort(8888),
		config.WithWorkPoolSize(1024),
		config.WithWorkerQueueCap(1024),
		config.WithMaxConn(500),
		config.WithTransportType(types.WebSocket),
	)
	if err != nil {
		logger.Fatal("failed to create server: %v", err)
	}

	// 注册消息
	p.RegisterMessage(MsgId_Login, func() any {
		return &LoginReq{}
	})

	// 注册路由
	r.Register(MsgId_Login, func(s types.Agent, msg types.MessagePack) {
		req := msg.GetMsg().(*LoginReq)
		// 业务逻辑
		logger.Info("Received LoginReq: %v", req.Code)
		s.WriteMsg(MsgId_LoginRes, &LoginRes{Code: req.Code + 1})
	})

	logger.Info("Server starting...")
	// 启动服务器
	if err := server.Run(); err != nil {
		logger.Fatal("failed to run server: %v", err)
	}
}
