package errorn

import (
	"strings"
	"testing"

	"nova/errorn/coden"
)

// TestErrorStack 测试错误堆栈功能
func TestErrorStack(t *testing.T) {
	// 1. 创建一个普通错误（不包含堆栈）
	err := New(coden.CommonUnknown, "test error")
	if err.StackTrace() != "" {
		t.Errorf("expected empty stack trace, got: %s", err.StackTrace())
	}

	// 2. 生成带堆栈信息的错误
	errWithStack := err.WithCallerStack(0)
	if errWithStack.StackTrace() == "" {
		t.Errorf("expected stack trace, but got empty string")
	}

	// 3. 确保堆栈信息有一定的关键字，比如本测试文件名
	if !strings.Contains(errWithStack.StackTrace(), "error_test.go") {
		t.Errorf("stack trace does not contain expected file info: %s", errWithStack.StackTrace())
	}

	// 4. Clone 复制错误后，堆栈信息应当保留
	clonedErr := Clone(errWithStack)
	if clonedErr.StackTrace() != errWithStack.StackTrace() {
		t.Errorf("expected stack trace to be copied, but got different results")
	}
	t.Log(clonedErr)

	// 5. 测试 FromError 是否能正确转换带有堆栈的错误
	fromErr := FromError(err)
	if fromErr.StackTrace() != err.StackTrace() {
		t.Errorf("expecting to convert an error with stack, but getting different results")
	}
	t.Log(fromErr)

	t.Log(err)
}

//
//func Test_ErrW(t *testing.T) {
//	err := errors.New("err msg")
//	err = NewRedis(err)
//	err2 := fmt.Errorf("err2; %w", err)
//	err3 := fmt.Errorf("err3; %w", err2)
//	logz.Info().Err(err).Msg("1")
//	logz.Info().Err(err2).Msg("2")
//	logz.Info().Err(err3).Msg("3")
//
//	errO2 := errors.Unwrap(err3)
//	logz.Info().Err(errO2).Msg("O2")
//}
