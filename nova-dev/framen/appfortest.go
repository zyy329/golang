package framen

import (
	"io"
	"os"
	"testing"

	"nova/logz"
	"nova/logz/zerologutil"

	"github.com/bytedance/mockey"
	"github.com/rs/zerolog"
)

var (
	testSpace    = []byte("    ")
	testSpaceNum = len(testSpace)
)

func TestMainRun(m *testing.M, components ...IComponent) {
	x := NewAppContext()

	// init
	if x.InitComponent(components...) != nil {
		return
	}

	// start
	if x.Start() != nil {
		return
	}

	// 调整logz日志输出格式, 使其更适合单元测试
	logz.SetLogger(TestLog())

	// 执行所有测试函数
	exitCode := m.Run()

	// 执行清理操作
	x.Showdown(0)

	// 注销 bytedance/mockey 注入
	mockey.UnPatchAll()

	// 退出测试，返回测试结果的退出码
	os.Exit(exitCode)
}

// TestLog 供单元测试使用的 log输出对象
func TestLog() (*zerolog.Logger, io.Closer) {
	writer := zerolog.ConsoleWriter{Out: newTestWriter(), NoColor: false, TimeFormat: zerolog.TimeFieldFormat}
	l := zerolog.New(writer).With().
		CallerWithSkipFrameCount(zerologutil.GetCallerSkip()).
		Logger().
		Level(zerolog.DebugLevel)

	return &l, writer
}

type testWriter struct {
	osWriter io.Writer
}

func newTestWriter() io.Writer {
	return &testWriter{osWriter: os.Stdout}
}

func (x *testWriter) Write(p []byte) (n int, err error) {
	// 比普通日志增加空格让格式更适合单元测试输出
	tp := make([]byte, 0, testSpaceNum+len(p))
	tp = append(tp, testSpace...)
	tp = append(tp, p...)
	n, err = x.osWriter.Write(tp)
	if err == nil {
		n -= testSpaceNum
	} else {
		// 如果写入失败，确保返回的 n 为 0
		n = 0
	}
	return n, err
}
