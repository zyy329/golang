package gneto

import (
	"fmt"

	"gneto/internal/core"
	"gneto/internal/logger"
	"gneto/pkg/config"
	"gneto/pkg/processor"
	"gneto/types"
)

// Server 游戏服务器实例
type Server struct {
	engine *core.Engine
}

// NewServer 创建服务器
func NewServer() *Server {
	return &Server{}
}

func (x *Server) Init(router types.RouterFunc, processor *processor.Processor, opts ...config.Option) error {
	options := config.NewDefaultOptions()
	for _, o := range opts {
		o(options)
	}

	if options.ByteOrder != processor.GetByteOrder() {
		return fmt.Errorf("server.byteOrder != processor.byteOrder; (%v/%v)", options.ByteOrder, processor.GetByteOrder())
	}

	x.engine = core.NewEngine(options, processor, router)

	return nil
}

// Run 启动服务器
func (x *Server) Run() error {
	return x.engine.Run()
}

// Stop 停止服务
func (x *Server) Stop() error {
	return x.engine.Stop()
}

// DispatchTo 将指定逻辑分发到指定工作协程进行处理
func (x *Server) DispatchTo(d types.DispatchAble, f func()) {
	x.engine.GetDispatcher().Dispatch(
		core.NewFuncTask(d, f))
}

//goland:noinspection GoUnusedExportedFunction
func SetLogger(log types.Logger) {
	logger.Init(log)
}
