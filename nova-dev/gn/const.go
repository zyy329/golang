package gn

import "time"

//goland:noinspection GoUnusedConst
const (
	Invalid = -1
)

/* Source 服务类型 -------------------------- */

type Source = string

//goland:noinspection GoUnusedConst
const (
	SourceGame       Source = "Game"
	SourceWeb        Source = "Web"
	SourceGameCenter Source = "GameCenter"
	SourceData       Source = "Data"
)

/* Source 服务类型 -------------------------- */

/* 时间常量定义 -------------------------- */
//goland:noinspection GoUnusedConst
const (
	SecMs    = 1000          // 毫秒数 -- 每秒
	MinuteMs = 60 * SecMs    // 毫秒数 -- 每分钟
	HourMs   = 60 * MinuteMs // 毫秒数 -- 每小时
	DayMs    = 24 * HourMs   // 毫秒数 -- 每天时
	WeekMs   = 7 * DayMs     // 毫秒数 -- 每周时
)

//goland:noinspection GoUnusedConst
const (
	MinuteSec = 60             // 秒数 -- 每分钟
	HourSec   = 60 * MinuteSec // 秒数 -- 每小时
	DaySec    = 24 * HourSec   // 秒数 -- 每天时
	WeekSec   = 7 * DaySec     // 秒数 -- 每周时
)

//goland:noinspection GoUnusedConst
const (
	TimeDay  = 24 * time.Hour
	TimeWeek = 7 * TimeDay
)

//goland:noinspection GoUnusedConst
const (
	MiniDateOnly = "20060102"
)

/* 时间常量定义 -------------------------- */

/* 日志常量定义 -------------------------- */
//goland:noinspection GoUnusedConst
const (
	SpanKey      = "span"
	TraceKey     = "trace"
	TimestampKey = "time_trace"
	LogIdKey     = "log_id"
)

/* 钱包类型定义 -------------------------- */

type WalletType = int32

//goland:noinspection GoUnusedConst
const (
	WalletTypeSingle   WalletType = 1 // 单一钱包
	WalletTypeTransfer WalletType = 2 // 转账钱包
)

/* 运营商类型定义 -------------------------- */

type OperatorType = int32

//goland:noinspection GoUnusedConst
const (
	OperatorTypeProduction OperatorType = 1 // 正式运营商
	OperatorTypeTest       OperatorType = 2 // 测试运营商
	OperatorTypeGray       OperatorType = 3 // 灰度运营商
)

//goland:noinspection GoUnusedConst
const (
	_        = iota             // 忽略第 0 个值
	KB int64 = 1 << (10 * iota) // 1 << 10
	MB                          // 1 << 20
	GB                          // 1 << 30
	TB                          // 1 << 40
	PB                          // 1 << 50
	EB                          // 1 << 60
)
