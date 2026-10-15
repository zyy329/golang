package internal

import (
	"fmt"

	"nova/logz"
	"nova/logz/logzlib"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func namingClient(url, namespace string) (client naming_client.INamingClient) {
	clientConf := createClientConf(namespace)
	var err error
	if client, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConf,
		ServerConfigs: createServerConf(url),
	}); err != nil {
		logz.Fatal().Msgf("nacos new namingClient err:%v", err)
		return nil
	}

	// 重设日志输出
	logger.SetLogger(logzlib.GetNacos(clientConf.LogLevel))

	return client
}

// Register 向 nacos 注册当前服务;
func Register(url, namespace string, p *vo.RegisterInstanceParam) (client naming_client.INamingClient) {
	client = namingClient(url, namespace)
	if client == nil {
		return
	}

	if success, err := client.RegisterInstance(*p); !success || err != nil {
		logz.Fatal().Msgf("nacos Register %v err:%v,%v", info(p), success, err)
		return nil
	}

	logz.Info().Msgf("nacos Register Success: %v", info(p))
	return
}

// UnRegister 主动注销当前服务
func UnRegister(client naming_client.INamingClient, p *vo.RegisterInstanceParam) {
	if client == nil {
		return
	}

	if success, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          p.Ip,
		Port:        p.Port,
		ServiceName: p.ServiceName,
		Ephemeral:   p.Ephemeral,
		Cluster:     p.ClusterName,
		GroupName:   p.GroupName,
	}); err != nil {
		logz.Fatal().Msgf("nacos UnRegister %v; err:%v,%v", info(p), success, err)
	} else {
		logz.Info().Msgf("nacos UnRegister Success: %v", info(p))
	}
	// 关闭连接
	client.CloseClient()
}

func Update(client naming_client.INamingClient, p *vo.RegisterInstanceParam) {
	if client == nil {
		return
	}

	if success, err := client.UpdateInstance(vo.UpdateInstanceParam{
		Ip:          p.Ip,
		Port:        p.Port,
		Weight:      p.Weight,
		Enable:      p.Enable,
		Healthy:     p.Healthy,
		Metadata:    p.Metadata,
		ClusterName: p.ClusterName,
		ServiceName: p.ServiceName,
		GroupName:   p.GroupName,
		Ephemeral:   p.Ephemeral,
	}); err != nil {
		logz.Info().Msgf("nacos Update %v, err:%v, %v", success, err, info(p))
	} else {
		logz.Info().Msgf("nacos Update Success: %v", info(p))
	}
}

func info(p *vo.RegisterInstanceParam) string {
	return fmt.Sprintf("[%v|%v:%v]--%v", p.ServiceName, p.Ip, p.Port, p.Metadata)
}
