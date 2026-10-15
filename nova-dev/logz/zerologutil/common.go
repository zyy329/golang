package zerologutil

const BaseCallerSkip = 2

// GetCallerSkip 获得跳过的调用栈层级
func GetCallerSkip() int {
	return BaseCallerSkip
}

type CloseEmpty struct {
}

func (x *CloseEmpty) Close() error {
	return nil
}
