package errorn

import (
	"fmt"

	"nova/errorn/coden"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	SlsLock = NewC(coden.DbSlsLockFailed)
)

//goland:noinspection GoUnusedExportedFunction
func NewUnknown(format string, v ...interface{}) *Error {
	return New(coden.CommonUnknown, sprintf(format, v...))
}

//goland:noinspection GoUnusedExportedFunction
func NewParam(format string, v ...interface{}) *Error {
	return New(coden.CommonParam, sprintf(format, v...))
}

//goland:noinspection GoUnusedExportedFunction
func NewMysql(err error) *Error {
	return New(coden.DbMysqlError, errorError(err))
}

//goland:noinspection GoUnusedExportedFunction
func NewMongo(err error) *Error {
	return New(coden.DbMongoError, errorError(err))
}

func NewRedis(err error) *Error {
	return New(coden.DbRedisError, errorError(err))
}

func NewSls(err error) *Error {
	return New(coden.DbSlsError, errorError(err))
}

func sprintf(format string, v ...interface{}) string {
	if len(v) == 0 {
		return format
	}

	cheatFunc := func(format string) string {
		return format
	}
	return fmt.Sprintf(cheatFunc(format), v...)
}

func errorError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
