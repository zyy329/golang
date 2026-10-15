package redisn

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	"nova/gn"
	"nova/gn/confkey"
	"nova/logz"
	"nova/nacos"

	"github.com/redis/go-redis/v9"
)

var (
	DefaultDataId = "redis"
	DefaultName   = confkey.Common
)

type RedisBaseConf struct {
	Addr        string
	Password    string
	TLSEnable   string
	PoolSize    int
	ReadTimeout string
}

func LoadRedisOption(name string) (*RedisBaseConf, error) {
	cf, err := loadAllConf()
	if err != nil {
		return nil, err
	}
	return cf[name], nil
}

func loadAllConf() (map[string]*RedisBaseConf, error) {
	content, err := nacos.LoadConf(DefaultDataId)
	if err != nil {
		return nil, err
	}
	cf := make(map[string]*RedisBaseConf)
	err = json.Unmarshal([]byte(content), &cf)
	if err != nil {
		return nil, err
	}
	return cf, nil
}

func fixConf(cf *RedisBaseConf) *redis.Options {
	opt := &redis.Options{
		Addr:        cf.Addr,
		Password:    cf.Password,
		DB:          0,
		PoolSize:    cf.PoolSize,
		ReadTimeout: time.Second * 30,
	}
	if cf.TLSEnable == "true" {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	if opt.PoolSize <= 0 {
		opt.PoolSize = 100
	}
	if cf.ReadTimeout != "" {
		duration, err := time.ParseDuration(cf.ReadTimeout)
		if err != nil {
			logz.Fatal().Msgf("Error parsing read timeout: %v", err)
		} else {
			opt.ReadTimeout = duration
		}
	}
	logz.Info().Msgf("redis addr %s", opt.Addr)
	return opt
}

func NewClient(nm ...string) *redis.Client {
	var name string
	if len(nm) == 0 {
		name = DefaultName
	} else {
		name = nm[0]
	}
	cf, err := LoadRedisOption(name)
	if err != nil {
		logz.Fatal().Msgf("failed redis init %v", err)
	}
	rdb := redis.NewClient(fixConf(cf))
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logz.Fatal().Msgf("failed connect redis %+v", err)
	}
	return rdb
}

//goland:noinspection GoUnusedExportedFunction
func Key(key string) string {
	return gn.SourceId + ":" + key
}
