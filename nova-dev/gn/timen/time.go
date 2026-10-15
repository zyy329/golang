package timen

import (
	"fmt"
	"time"

	"nova/gn"
)

func NowTs() int64 {
	return time.Now().Unix()
}

//goland:noinspection GoUnusedExportedFunction
func NowTsms() int64 {
	return time.Now().UnixMilli()
}

func TsToTime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func TsmsToTime(tsms int64) time.Time {
	return time.UnixMilli(tsms)
}

func TimeToTs(t time.Time) int64 {
	return t.Unix()
}

//goland:noinspection GoUnusedExportedFunction
func TimeToTsms(t time.Time) int64 {
	return t.UnixMilli()
}

func MsToSec(ms int64) int64 {
	return ms / gn.SecMs
}

//goland:noinspection GoUnusedExportedFunction
func SecToMs(sec int64) int64 {
	return sec * gn.SecMs
}

// ToDayZero 计算目标时间对应的 当天零点时间戳; 当地时区
func ToDayZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// ToMonthZero 计算目标时间对应的 当月一号零点时间戳; 当地时区
func ToMonthZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

//goland:noinspection GoUnusedExportedFunction
func NextDay(ts int64) int64 {
	_, next := GetCycleTsByTs(CycleTypeDay, ts)
	return next
}

//goland:noinspection GoUnusedExportedFunction
func NextMonth(ts int64) int64 {
	_, next := GetCycleTsByTs(CycleTypeMonth, ts)
	return next
}

// DiffDay 获取2天的日期差 精确到天
//
//goland:noinspection GoUnusedExportedFunction
func DiffDay(start time.Time, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}

// GetLastDateOfMonth 获取本月最后一天日期
//
//goland:noinspection GoUnusedExportedFunction
func GetLastDateOfMonth(date time.Time) time.Time {
	return ToDayZero(date).AddDate(0, 1, 0).AddDate(0, 0, -1)
}

//goland:noinspection GoUnusedExportedFunction
func PreviousDate(date time.Time) time.Time {
	return date.AddDate(0, 0, -1)
}

//goland:noinspection GoUnusedExportedFunction
func NextDate(date time.Time) time.Time {
	return date.AddDate(0, 0, 1)
}

//goland:noinspection GoUnusedExportedFunction
func GetLocalHourOffset() int {
	_, offset := time.Now().Zone()
	return offset / 3600
}

func GetLocation(zone float64) *time.Location {
	return time.FixedZone(fmt.Sprintf("UTC%+.2f", zone), int(zone*gn.HourSec))
}

func SetZone(zone float64) {
	time.Local = GetLocation(zone)
}
