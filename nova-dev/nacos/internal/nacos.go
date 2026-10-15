package internal

import (
	"strconv"
	"strings"

	"nova/logz"
	"nova/logz/logzlib"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func newClient(url, namespace string) (client config_client.IConfigClient, err error) {
	clientConf := createClientConf(namespace)
	if client, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConf,
		ServerConfigs: createServerConf(url),
	}); err != nil {
		logz.Fatal().Msgf("nacos NewConfigClient err:%v", err)
	}

	// 重设日志输出
	logger.SetLogger(logzlib.GetNacos(clientConf.LogLevel))

	return client, err
}

func createClientConf(namespace string) (conf *constant.ClientConfig) {
	conf = constant.NewClientConfig(
		constant.WithNamespaceId(func() string {
			return namespace
		}()),
		constant.WithTimeoutMs(5000),
		constant.WithCacheDir("./_cache_nacos"),
		constant.WithLogLevel("warn"),
		constant.WithLogRollingConfig(&constant.ClientLogRollingConfig{
			MaxBackups: 3,
			LocalTime:  true,
		}),
	)
	conf.AppendToStdout = true
	return
}

func createServerConf(confUrls string) (conf []constant.ServerConfig) {
	if len(confUrls) == 0 {
		logz.Fatal().Msgf("API_NACOS is empty")
	}

	var urlS []string
	if urlS = strings.Split(confUrls, ","); len(urlS) == 0 {
		logz.Fatal().Msgf("API_NACOS Split err:%v", confUrls)
	}

	for idx, url := range urlS {
		if itemS := strings.Split(url, ":"); len(itemS) != 2 {
			logz.Fatal().Msgf("API_NACOS ip port Split err:%v, idx:%v", confUrls, idx)
		} else {
			if port, err := strconv.Atoi(itemS[1]); err != nil {
				logz.Fatal().Msgf("API_NACOS ip port Split err:%v, idx:%v", confUrls, idx)
			} else {
				conf = append(conf, constant.ServerConfig{IpAddr: itemS[0], Port: uint64(port)})
			}
		}
	}

	return
}

func getContent(client config_client.IConfigClient, dataId, group string) (content string, err error) {
	return client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}
