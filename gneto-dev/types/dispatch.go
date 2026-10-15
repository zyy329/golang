package types

type DispatchAble interface {
	DispatchKey() int64
	DispatchFailed() // 分发失败会触发
}
