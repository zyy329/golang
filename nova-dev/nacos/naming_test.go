package nacos

import (
	"testing"
	"time"

	"nova/framen"
	"nova/gn/conv"
	"nova/gn/timen"

	"github.com/bytedance/mockey"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// 测试启动, 做统一的初始化化工作, 同包下每个Test都会自动调用该函数
func TestMain(m *testing.M) {
	framen.TestMainRun(m)
}

func TestNaming_UpdateMeta(t *testing.T) {
	mockey.PatchConvey("Test Nacos Naming Update", t, func() {
		key := "num"
		meta := map[string]string{
			"gameId": "9999",
			"node":   "testNode",
			"ts":     conv.ToStr(timen.NowTs()),
			key:      "0",
		}

		n := NewNaming()
		n.Register(vo.RegisterInstanceParam{
			Ip:          "127.0.0.1",
			Port:        6666,
			Weight:      1,
			Enable:      true,
			Healthy:     true,
			Metadata:    meta,
			ServiceName: "TestUp",
			Ephemeral:   true,
		})

		for i := 1; i <= 10; i++ {
			time.Sleep(5 * time.Second)
			meta[key] = conv.ToStr(i)
			n.UpdateMeta(meta)
		}

		n.UnRegister()
	})
}
