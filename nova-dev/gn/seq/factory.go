package seq

import (
	"math/rand/v2"
	"sync"
	"time"

	"nova/gn"
	"nova/gn/envn"
)

const (
	yearCycle  = 2
	tsMusk     = 100000000000
	moduleMusk = 10000000
	podMusk    = 100000
	tolerance  = 60              // 各节点可能时间有误差
	maxExpire  = gn.DaySec * 364 // 保存时间 < 1年

	maxIncId = 99999 // 1s只能产生5位数以内
	startIdB = 10000

	yearBaseTs0 = 10000000 // 目前分成两级没有什么用，后面可能用来区分奇偶年
	yearBaseTs1 = 50000000
)

type Factory struct {
	lock sync.Mutex

	module      int32
	lastTs      int64
	yearBaseTsS [yearCycle]int64

	incId   int64
	startId int64
}

// NewFactory 不同的功能、模块，使用独立的Factory，增加并发
func NewFactory(moduleId int32) *Factory {
	x := &Factory{
		module:      moduleId,
		yearBaseTsS: [yearCycle]int64{yearBaseTs0, yearBaseTs1},
		startId: func() (start int64) {
			if envn.IsPodB {
				return startIdB
			}
			return
		}(),
	}

	return x
}

/*
	产生全局唯一序列号

限制：
1. 数据最多保存6个月
2. NodeId最大为99
3. ModuleId需要初始化，最大为9999
*/
func (x *Factory) Id() (id int64, t time.Time) {
	var incId int64
	t, incId = x.timeAndIncId()

	id = func() int64 {
		year := t.Year()
		return int64(t.Sub(time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)).Seconds()) +
			x.yearBaseTsS[year%yearCycle]
	}()*tsMusk +
		int64(x.module)*moduleMusk +
		int64(envn.PodId)*podMusk +
		incId

	return
}

func (x *Factory) timeAndIncId() (t time.Time, incId int64) {
	x.lock.Lock()
	defer x.lock.Unlock()

	t = time.Now()
	if ts := t.Unix(); ts != x.lastTs {
		x.lastTs = ts
		x.resetIncId()
	} else if x.incId < maxIncId {
		x.incId++
	} else {
		time.Sleep(time.Until(t.Truncate(time.Second).Add(time.Second)))
		t = time.Now()
		x.lastTs = t.Unix()
		x.resetIncId()
	}
	incId = x.incId
	return
}

// 后面不够的话，可以改成一个往上涨，一个往下降，这样并发可以X2
func (x *Factory) resetIncId() {
	x.incId = x.startId + rand.Int64N(5)
}
