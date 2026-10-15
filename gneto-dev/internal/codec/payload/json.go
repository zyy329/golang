package payload

import (
	"encoding/json"
	"reflect"
)

// IResettable 允许结构体自实现高性能 Reset
type IResettable interface {
	Reset()
}
type JsonPayload struct {
	*MessagePool
}

func NewJsonPayload(pool *MessagePool) *JsonPayload {
	p := &JsonPayload{MessagePool: pool}
	return p
}

func (x *JsonPayload) PushMessage(msgId uint16, msg any) {
	if msg == nil {
		return
	}

	// 优先使用接口定义的 Reset，性能最高
	if r, ok := msg.(IResettable); ok {
		r.Reset()
	} else {
		// 使用反射将结构体置回零值，方便尽早暴露业务逻辑中不合理的持续持有; 后续如发现性能问题, 建议业务层提供 IResettable 实现来优化
		v := reflect.ValueOf(msg)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
			if v.CanSet() {
				v.Set(reflect.Zero(v.Type()))
			}
		}
	}
	x.pushToPool(msgId, msg)
}
func (x *JsonPayload) Marshal(msg any) ([]byte, error) {
	return json.Marshal(msg)
}
func (x *JsonPayload) Unmarshal(data []byte, msg any) error {
	return json.Unmarshal(data, msg)
}
