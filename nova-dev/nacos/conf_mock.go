package nacos

import (
	"nova/gn/envn"
	"nova/nacos/internal"
)

var (
	mockConf = make(map[string]map[string]string)
)

//go:cover exclude 该函数不计入单元测试覆盖率计算
func AddMockConf(group string, dataId string, confStr string) {
	if mockConf[group] == nil {
		mockConf[group] = make(map[string]string)
	}
	mockConf[group][dataId] = confStr
}

// MockLoadConf 返回 AddMockConf 设定的数据,而不从远程nacos中去读取;
//
//go:cover exclude 该函数不计入单元测试覆盖率计算
func MockLoadConf() {
	if !envn.MockSign {
		return
	}

	internal.LoadFunc = func(url, namespace, dataId, group string) (string, error) {
		if mockConf[group] != nil {
			return mockConf[group][dataId], nil
		} else {
			return "", nil
		}
	}
}
