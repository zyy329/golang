package seq

import (
	"time"

	"nova/gn/timen"
)

func IdToTs(id int64) (ts int64) {
	now := time.Now()
	nowTs := now.Unix()
	year := now.Year()

	yearSec := id / tsMusk
	if yearSec >= yearBaseTs1 {
		yearSec -= yearBaseTs1
	} else {
		yearSec -= yearBaseTs0
	}

	ts = time.Date(year, 1, 1, 0, 0, 0, 0, time.Local).Unix() + yearSec
	if diff := ts - nowTs; diff > tolerance {
		ts = time.Date(year-1, 1, 1, 0, 0, 0, 0, time.Local).Unix() + yearSec
	} else if -diff > maxExpire {
		ts = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local).Unix() + yearSec
	}
	return
}

func IdToTime(id int64) time.Time {
	return timen.TsToTime(IdToTs(id))
}

func IdToModuleId(id int64) int {
	return int(id % tsMusk / moduleMusk)
}
