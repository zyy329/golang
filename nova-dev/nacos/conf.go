package nacos

import (
	"nova/nacos/internal"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
)

// nacos 配置中心 支持

const (
	DefaultGroup = "DEFAULT_GROUP"
)

var (
	// 当前环境变量(url, namespace)关联的 nacos client
	client config_client.IConfigClient
)

// LoadConf 加载获得 nacos 中的配置数据
// @param	param:   nacos 连接所需参数; 如果传入为 nil, 会尝试从 环境变量中去获取
// @param	group:   nacos group信息; 不填默认值为 DEFAULT_GROUP; 填写多个, 只会有第一个参数有效
// @return  content: 原始配置数据, 不进行 反序列化解析
func LoadConf(dataId string, group ...string) (string, error) {
	var g string
	if len(group) > 0 {
		g = group[0]
	} else {
		g = DefaultGroup
	}

	param := LoadByEnv()

	return internal.LoadFunc(param.Url, param.Namespace, dataId, g)
}

//goland:noinspection GoUnusedExportedFunction
func LoadConfByConnParam(param *ConnParam, dataId string, group ...string) (string, error) {
	var g string
	if len(group) > 0 {
		g = group[0]
	} else {
		g = DefaultGroup
	}

	return internal.LoadFunc(param.Url, param.Namespace, dataId, g)
}

//goland:noinspection GoUnusedExportedFunction
func GetClient() (config_client.IConfigClient, error) {
	if client == nil {
		param := LoadByEnv()
		if c, err := internal.GetClient(param.Url, param.Namespace); err != nil {
			return nil, err
		} else {
			client = c
		}
	}
	return client, nil
}

//goland:noinspection GoUnusedExportedFunction
func GetClientByConnParam(param *ConnParam) (config_client.IConfigClient, error) {
	return internal.GetClient(param.Url, param.Namespace)
}
