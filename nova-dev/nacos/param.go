package nacos

import (
	"os"

	"nova/gn/envn"
)

// nacos 连接参数

type ConnParam struct {
	// nacos 地址
	Url string
	// nacos 命名空间
	Namespace string
}

func LoadByEnv() *ConnParam {
	return &ConnParam{
		Url: os.Getenv("API_NACOS"),
		Namespace: func() string {
			if envn.Env != envn.Prod {
				if namespace := os.Getenv("API_NACOS_NAMESPACE"); namespace != "" {
					return namespace
				}
			}
			return envn.EnvName(func() (project string) {
				if project = os.Getenv("API_PROJECT"); project == "" {
					project = "api"
				}
				return
			}())
		}(),
	}
}
