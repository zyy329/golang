package etcdn

import (
	"encoding/json"

	"nova/nacos"

	"github.com/pkg/errors"
)

const (
	Etcd = "etcd"
)

type EtcdConf struct {
	Hosts []string
	User  string
	Pass  string
}

func LoadEtcd(name string) (*EtcdConf, error) {
	content, err := nacos.LoadConf(Etcd)
	if err != nil {
		return nil, errors.WithMessage(err, "nacos load etcd config")
	}
	if content == "" {
		return nil, nil
	}
	var cfg map[string]*EtcdConf
	if err = json.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, errors.WithMessage(err, "config content parse")
	}
	return cfg[name], nil
}
