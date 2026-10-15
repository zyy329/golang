package nacos

import (
	"sync/atomic"

	"nova/nacos/internal"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// nacos 服务发现 支持

type Naming struct {
	client    naming_client.INamingClient
	sign      int32 // 标记, 用于确保 注册 和 删除 整个对象生命周期中都只会被调用一次
	param     *vo.RegisterInstanceParam
	connParam *ConnParam
}

func NewNaming() *Naming {
	return NewNamingByConnParam(LoadByEnv())
}

func NewNamingByConnParam(param *ConnParam) *Naming {
	return &Naming{connParam: param}
}

// Register 向 nacos 注册当前服务;
func (x *Naming) Register(p vo.RegisterInstanceParam) {
	if !atomic.CompareAndSwapInt32(&x.sign, 0, 1) {
		// 避免重复操作
		return
	}

	x.param = &p
	x.client = internal.Register(x.connParam.Url, x.connParam.Namespace, x.param)
}

// UnRegister 主动注销当前服务
func (x *Naming) UnRegister() {
	// 特意将 new 设置为2, 使得 注册 和 删除 整个进程生命周期中都只会被调用一次
	if !atomic.CompareAndSwapInt32(&x.sign, 1, 2) {
		// 避免重复操作
		return
	}

	internal.UnRegister(x.client, x.param)
}

// UpdateMeta 更新元数据
func (x *Naming) UpdateMeta(metadata map[string]string) {
	x.param.Metadata = metadata
	if atomic.LoadInt32(&x.sign) == 1 {
		internal.Update(x.client, x.param)
	}
}
