package errorn

import "runtime"

const (
	DefaultStackLimit = 5
)

type StackFrame = uintptr

// GetStackS 生成数组形式的堆栈信息; skip, 跳过的堆栈层数 limit, 最大记录的堆栈层数
func GetStackS(skip, limit int) []StackFrame {
	pcs := make([]StackFrame, limit)
	n := runtime.Callers(skip+2, pcs) // +2 跳过 GetStackS 和 runtime.Callers 本身
	if n == 0 {
		return nil
	} else {
		return pcs[:n]
	}
}
