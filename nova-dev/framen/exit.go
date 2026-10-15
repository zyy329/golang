package framen

import (
	"os"
	"os/signal"
	"syscall"

	"nova/logz"
)

// WaitExitSig 阻塞等待退出信号
func WaitExitSig() {
	c := make(chan os.Signal, 1)
	// syscall.SIGTERM(15), docker/k8s 关闭前 发送的的是 syscall.SIGTERM 信号(信号后,默认等待30秒优雅关闭时间才正式关闭)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM)
	sig := <-c
	logz.Info().Msgf("Get Exit signal: %v", sig)
}

// AsyncWaitExitSig 异步等待退出信号; 接到信号后会触发回调
//
//goland:noinspection GoUnusedExportedFunction
func AsyncWaitExitSig(callback func()) {
	go func() {
		WaitExitSig()
		callback()
	}()
}
