package conv

import (
	"fmt"
	"strconv"
)

//goland:noinspection GoUnusedExportedFunction
func AtoiByDefault[T int | int32 | int64](str string, def T) T {
	if n, err := strconv.ParseInt(str, 10, 64); err != nil {
		return def
	} else {
		return T(n)
	}
}

// ToStr 基础类型转换为Str
// 目前 支持的类型
//
//	int..., uint..., bool, string
//	使用断言, 效率更优
//
// 没有支持的类型, 使用 fmt.Sprintf, 效率低一点
func ToStr(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", val)
	}
}

//goland:noinspection GoUnusedExportedFunction
func StrToInt[T int | int32 | int64](str string) T {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	} else {
		return T(i)
	}
}

type BasicType interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		bool | string
}

// IsZeroValue 基础类型专用版：检查输入是否为其类型的默认值
//
//goland:noinspection GoUnusedExportedFunction
func IsZeroValue[T BasicType | *any](v T) bool {
	var zero T
	return v == zero
}
