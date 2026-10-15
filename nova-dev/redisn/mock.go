package redisn

import (
	"testing"

	"nova/gn/envn"

	"github.com/alicebob/miniredis/v2"
	"github.com/bytedance/mockey"
)

// MockRedis 单元测试时, 注入redis客户端;
// 利用 miniredis, 启动一个本地内存redis; 将对redis的访问重定向到该本地redis; 避免对外部redis的依赖
// miniredis 官网地址:https://github.com/alicebob/miniredis
//
//go:cover exclude 该函数不计入单元测试覆盖率计算
func MockRedis(t *testing.T) *miniredis.Miniredis {
	// 启动一个模拟redis
	mr := miniredis.RunT(t)

	if !envn.MockSign {
		// 没有进行 mock, 实际miniredis 无效果, 不影响测试
		return mr
	}

	// mock 的方式重定向启动的Redis客户端连接到该本地模拟实例
	// redisn.NewClient 方法创建redis客户端时, 会调用loadRedisOption函数获取连接配置
	// 这里通过注入loadRedisOption函数, 修改返回配置中的地址重定向到模拟redis, 进而到达 mock 替代的作用
	origin := LoadRedisOption
	mockey.Mock(LoadRedisOption).
		To(func(name string) (*RedisBaseConf, error) {
			// 调用原始配置函数
			if cfg, err := origin(name); err != nil {
				return cfg, err
			} else {
				// 修改redis 地址, 重定向到本地模拟实例
				cfg.Addr = mr.Addr()
				return cfg, nil
			}
		}).Origin(&origin).Build()

	return mr
}
