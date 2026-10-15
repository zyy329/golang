package zerologutil

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"nova/errorn"
	"nova/gn/envn"
	"nova/gn/strn"
	"nova/logz/logzdef"

	"github.com/rs/zerolog"
)

const (
	// 堆栈打印时, 打印的文件路径层数
	stackPathLevel   = 3
	stackPathCellNum = stackPathLevel + 1
	LogTimeFormat    = "01-02 15:04:05.000000"
)

func init() {
	// 时间输出格式
	zerolog.TimeFieldFormat = LogTimeFormat

	// 默认最简 调用栈输出信息
	ChangeCallerMarshalFunc(0)

	zerolog.MessageFieldName = "msg"
	zerolog.ErrorStackMarshaler = MarshalStack
}

//goland:noinspection GoUnusedExportedFunction
func DiscardLog() (*zerolog.Logger, io.Closer) {
	l := zerolog.New(nil).Level(zerolog.Disabled)
	return &l, &CloseEmpty{}
}

// DefaultLog 默认日志输出, 用于还未初始化好之前的日志输出
func DefaultLog() (*zerolog.Logger, io.Closer) {
	// 本地调试环境才开启色彩; 云上环境不开
	noColor := envn.Env != envn.Local
	writer := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: noColor, TimeFormat: zerolog.TimeFieldFormat}
	l := zerolog.New(writer).With().
		CallerWithSkipFrameCount(GetCallerSkip()).
		Logger().
		Level(zerolog.InfoLevel)

	return &l, writer
}

// ChildByLevel 创建主log 对象的 指定日志级别的子副本;
func ChildByLevel(log *zerolog.Logger, level logzdef.LogLevel) *zerolog.Logger {
	l := log.Level(logzdef.ToZeroLevel(level))
	return &l
}

// ChangeCallerMarshalFunc 修改调用栈输出信息 显示多少级路径
func ChangeCallerMarshalFunc(pathLevel int) {
	if pathLevel <= 0 {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
	} else {
		cellNum := pathLevel + 1
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			cellS := strings.Split(file, string(filepath.Separator))
			if len(cellS) > cellNum {
				file = filepath.Join(cellS[len(cellS)-cellNum:]...)
			}
			return file + ":" + strconv.Itoa(line)
		}
	}
}

func MarshalStack(err error) interface{} {
	var e *errorn.Error
	// 只处理来自 errorn.Error 的堆栈
	if !errors.As(err, &e) {
		return nil
	} else {
		return &LogStack{e.StackS()}
	}
}

type LogStack struct {
	StackS []errorn.StackFrame
}

func (x *LogStack) MarshalZerologObject(l *zerolog.Event) {
	LogPrintStack(l, x.StackS)
}

func LogPrintStack(l *zerolog.Event, stackS []errorn.StackFrame) *zerolog.Event {
	stLen := len(stackS)
	if stLen == 0 {
		// 没有堆栈信息
		return l
	}

	// 生成堆栈信息
	frames := runtime.CallersFrames(stackS)
	for i := 0; i < stLen; i++ {
		frame, more := frames.Next()

		// 只生成简化的堆栈信息: 有限层数路径/文件名+行号
		cellS := strings.Split(frame.File, string(filepath.Separator))
		var file string
		if len(cellS) > stackPathCellNum {
			file = filepath.Join(cellS[len(cellS)-stackPathCellNum:]...)
		} else {
			file = frame.File
		}
		l.Str(fmt.Sprintf("%s%d", strn.Stack, i), file+":"+strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}
	return l
}
