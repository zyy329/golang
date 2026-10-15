package errorn

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"nova/errorn/coden"

	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code     int32             `json:"code"`
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	stack    []byte
	stackS   []StackFrame // 原始堆栈数组信息
}

func (x *Error) Error() string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "code=%d", x.Code)
	if x.Message != "" {
		builder.WriteString(", message=")
		builder.WriteString(x.Message)
	}
	if len(x.Metadata) > 0 {
		builder.WriteString(", metadata=[")
		for k, v := range x.Metadata {
			builder.WriteString("(")
			builder.WriteString(k)
			builder.WriteString(":")
			builder.WriteString(v)
			builder.WriteString(")")
		}
		builder.WriteString("]")
	}
	return builder.String()
}

func (x *Error) Is(err error) bool {
	if x == nil {
		return err == nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se.Code == x.Code
	}
	return false
}

func (x *Error) IsCode(code int32) bool {
	if x == nil {
		return false
	}
	return x.Code == code
}

func (x *Error) SetMessage(message string) *Error {
	x.Message = message
	return x
}

func (x *Error) SetMetadata(key string, value string) *Error {
	if x.Metadata == nil {
		x.Metadata = make(map[string]string)
	}
	x.Metadata[key] = value
	return x
}

func (x *Error) WithMetadata(errors map[string]string) *Error {
	x.Metadata = errors
	return x
}

// WithDebugStack 生成完整的 goroutine 堆栈
func (x *Error) WithDebugStack() *Error {
	x.stack = debug.Stack()
	return x
}

// WithStack 仅记录简要调用栈
func (x *Error) WithStack() *Error {
	return x.WithCallerStack(1)
}

// WithCallerStack 仅记录简要调用栈
func (x *Error) WithCallerStack(skip int) *Error {
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip+2, pcs) // +2 跳过 WithCallerStack 和 runtime.Callers 本身
	if n == 0 {
		return x
	}
	var builder strings.Builder
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		_, _ = fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	x.stack = []byte(builder.String())
	return x
}

// StackTrace 返回错误的堆栈信息（字符串）
func (x *Error) StackTrace() string {
	return string(x.stack)
}

// Stack 返回错误的堆栈信息
func (x *Error) Stack() []byte {
	return x.stack
}

// WithCallerStackS 生成数组形式的堆栈信息; skip, 跳过的堆栈层数 limit, 最大记录的堆栈层数
func (x *Error) WithCallerStackS(skip, limit int) *Error {
	x.stackS = GetStackS(skip+1, limit) // +1 跳过 WithCallerStackS 本身
	return x
}
func (x *Error) WithStackS() *Error {
	return x.WithCallerStackS(1, DefaultStackLimit)
}

func (x *Error) StackS() []StackFrame {
	return x.stackS
}

// GRPCStatus 实现 GRPCStatus() 接口
func (x *Error) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(x.Code), x.Message).
		WithDetails(&errdetails.ErrorInfo{
			Metadata: x.Metadata,
		})
	return s
}

// New 创建错误
func New(code int32, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewC 创建错误（只有错误code，后续再设置描述）
func NewC(code int32) *Error {
	return &Error{
		Code: code,
	}
}

// Clone 深拷贝到新的 Error
func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	var stack []byte
	if len(err.stack) > 0 {
		stack = append([]byte(nil), err.stack...)
	}
	return &Error{
		Code:     err.Code,
		Message:  err.Message,
		Metadata: metadata,
		stack:    stack,
	}
}

// FromError 转换错误为 Error
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(coden.CommonUnknown, err.Error())
	}
	ret := New(
		int32(gs.Code()), gs.Message(),
	)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

// IsCode 判断错误是否为指定code
func IsCode(err error, code int32) bool {
	if err == nil {
		return false
	}
	return FromError(err).Code == code
}

// Is 判断错误是否为指定错误
//
//goland:noinspection GoUnusedExportedFunction
func Is(err, target error) bool {
	if target == nil {
		return err == nil
	}
	return IsCode(err, FromError(target).Code)
}

// ToServiceErr 检查错误类型
// 错误码为预期错误码体系内的可以直接抛出 如玩家不存在  参数错误等
// 不在错误码体系内的 如json序列化错误、空指针等  转为 定义的未知错误类型
// 建议返回错误码数字作为错误编码便于追查错误，但不推荐将数据库异常内部描述之类的内容不经截断直接明文返回到api的字段上，
// 如driver驱动内部代码定义的内容直接赋值到message中，但必须输出详细日志，建议输出堆栈信息便于追踪错误发生点
func ToServiceErr(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var e *Error
	if errors.As(err, &e) && e.Code >= coden.HttpOk {
		if e.Code >= coden.DbError && e.Code <= coden.DbSlsError {
			return New(e.Code, "db error"), true
		}
		return e, true
	} else {
		return New(coden.CommonInternal, "internal error"), false
	}
}
