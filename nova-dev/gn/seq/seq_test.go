package seq

import (
	"testing"

	"nova/gn/timen"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
)

func TestFactory(t *testing.T) {
	seqFactory := NewFactory(5117)
	for i := 0; i < 1; i++ {
		id, idTime := seqFactory.Id()
		ts := IdToTs(id)
		time := IdToTime(id)
		t.Logf("id:%v idTime:%v ts:%v  time:%v ", id, idTime, ts, time)
	}
}

func TestIdTime0(t *testing.T) {
	id := int64(7612977351170000000)
	ts := IdToTs(id)
	time := timen.TsToTime(ts)
	t.Logf("id:%v  ts:%v  time:%v", id, ts, time)
}

func TestIdTime1(t *testing.T) {
	id := int64(7189484374751170061)
	ts := IdToTs(id)
	time := timen.TsToTime(ts)
	t.Logf("id:%v  ts:%v  time:%v", id, ts, time)
}

func TestIdToModuleId(t *testing.T) {
	mockey.PatchConvey("SeqIdToModuleId", t, func() {
		mId := int32(9900)
		f := NewFactory(mId)
		id, _ := f.Id()
		moduleId := IdToModuleId(id)
		convey.So(moduleId, convey.ShouldEqual, mId)
	})
}
