package timen

import (
	"fmt"
	"time"
)

// CycleType 时间周期工具函数
type CycleType int32

//goland:noinspection GoUnusedConst
const (
	CycleTypeUnknown CycleType = 0
	CycleTypeDay     CycleType = 1
	CycleTypeMonth   CycleType = 3

	CycleStrDay   = "day"
	CycleStrMonth = "month"
)

//goland:noinspection GoUnusedExportedFunction
func GetCycleStr(cycleType CycleType) string {
	switch cycleType {
	case CycleTypeDay:
		return CycleStrDay
	case CycleTypeMonth:
		return CycleStrMonth
	default:
		panic(fmt.Sprintf("invalid cycle: %v", cycleType))
	}
}

// GetCycleTs 获得指定周期的起始和结束时间，传入不存在的周期会panic
// 区间关系 [start, end)
// 比如: 周期:天, 实际返回 当天0点 和 第二天0点; 需要当成左闭右开 [当天0点,第二天0点) 才是正确的一周期时间描述
func GetCycleTs(cycleType CycleType, t time.Time) (startTs, endTs int64) {
	switch cycleType {
	case CycleTypeDay:
		return TimeToTs(ToDayZero(t)), TimeToTs(ToDayZero(t.AddDate(0, 0, 1)))
	case CycleTypeMonth:
		return TimeToTs(ToMonthZero(t)), TimeToTs(ToMonthZero(t.AddDate(0, 1, 0)))
	default:
		panic(fmt.Sprintf("invalid cycle: %v", cycleType))
	}
}

func GetCycleTsByTs(cycleType CycleType, ts int64) (startTs, endTs int64) {
	return GetCycleTs(cycleType, TsToTime(ts))
}
