package internal

import (
	"sync"

	"nova/logz"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
)

var (
	clientM  = make(map[string]config_client.IConfigClient)
	lock     sync.Mutex
	LoadFunc = func(url, namespace, dataId, group string) (string, error) {
		return LoadConf(url, namespace, dataId, group)
	}
)

func GetClient(url, namespace string) (config_client.IConfigClient, error) {
	lock.Lock()
	defer lock.Unlock()

	key := namespace + "_" + url
	c, ok := clientM[key]
	if ok {
		return c, nil
	} else {
		client, err := newClient(url, namespace)
		if err != nil {
			return nil, err
		}
		clientM[key] = client

		//// 仍然会报错, 不关闭了
		//time.AfterFunc(10*time.Second, func() {
		//	lock.Lock()
		//	defer lock.Unlock()
		//
		//	// 延迟关闭; 减少重复建立连接, 避免 client connection is closing 报错
		//	clientM[key].CloseClient()
		//	delete(clientM, key)
		//	logz.Info().Msgf("delete")
		//})

		return client, nil
	}
}

func LoadConf(url, namespace, dataId, group string) (string, error) {
	client, err := GetClient(url, namespace)
	if err != nil {
		return "", err
	}

	content, err2 := getContent(client, dataId, group)
	if err2 != nil {
		logz.Fatal().Msgf("LoadConf; url: %s, namespace: %s, dataId: %s, group: %s; suc:%v, size:%v",
			url, namespace, dataId, group, err2, len(content))
	}
	logz.Info().Msgf("LoadConf; url: %s, namespace: %s, dataId: %s, group: %s;  size:%v",
		url, namespace, dataId, group, len(content))
	return content, err2
}
