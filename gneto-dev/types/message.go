package types

type MessagePack interface {
	// GetMsgId 消息Id
	GetMsgId() uint16
	// GetMsg 获取原始消息体
	// 注意:该消息对象不可持久化持有或传入其他协程,否则可能出现数据错误(因为直接调用逻辑结束后,该消息会被清理和还回重用池)
	GetMsg() any
	// AcquireMsg 获取原始消息,并转移所有权(MsgPack 内的 msg 变为 nil，且不回再被重用池回收); 业务使用完后交由GC回收
	// 用于 应对GetMsg的不可持有性, 当有特殊需求要持久化持有或传入其他协程需求时, 可使用该方法
	AcquireMsg() any
}
