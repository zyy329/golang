package payload

import (
	"errors"
	"fmt"
	"sync"

	"gneto/types"
)

var (
	ErrMessageNotRegister = errors.New("message not register")
)

// MessagePool 提供 payload 通用的对象池管理逻辑
type MessagePool struct {
	// msgPools 存储每个 msgId 对应的 sync.Pool, 减少频繁创建和GC
	msgPools    map[uint16]*sync.Pool
	msgFactoryS map[uint16]types.MessageCreator
}

func NewMessagePool() *MessagePool {
	return &MessagePool{
		msgPools:    make(map[uint16]*sync.Pool),
		msgFactoryS: make(map[uint16]types.MessageCreator),
	}
}

func (b *MessagePool) Register(msgID uint16, factory types.MessageCreator) {
	b.msgPools[msgID] = &sync.Pool{
		New: func() any { return factory() },
	}
	b.msgFactoryS[msgID] = factory
}

func (b *MessagePool) NewMessage(msgID uint16) (any, error) {
	factory, ok := b.msgFactoryS[msgID]
	if !ok {
		return nil, fmt.Errorf("%w; id: %d", ErrMessageNotRegister, msgID)
	}
	return factory(), nil
}

func (b *MessagePool) PopMessage(msgID uint16) (any, error) {
	pool, ok := b.msgPools[msgID]
	if !ok {
		return nil, fmt.Errorf("%w; id: %d", ErrMessageNotRegister, msgID)
	}
	return pool.Get(), nil
}

func (b *MessagePool) pushToPool(msgID uint16, msg any) {
	if pool, ok := b.msgPools[msgID]; ok {
		pool.Put(msg)
	}
}
