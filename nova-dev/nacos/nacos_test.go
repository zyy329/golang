package nacos

import (
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func TestConf(t *testing.T) {
	AddMockConf(DefaultGroup, "collect", "aaa")
	MockLoadConf()

	conf, _ := LoadConf("collect")

	t.Log(conf)
}

func TestNaming(t *testing.T) {
	n := NewNaming()
	//n := nacos.NewNamingByConnParam(&nacos.ConnParam{Url: "127.0.0.0:800", Namespace: "dev"})
	n.Register(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        5000,
		ServiceName: "game",
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true, // 临时模式, 服务器关闭后, 缺失心跳会自动注销
		Metadata:    map[string]string{"gameId": "800"},
	})

	time.Sleep(10 * time.Second)
	n.UnRegister()

	time.Sleep(10 * time.Second)
}
