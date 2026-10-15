package framen

import (
	"time"

	"nova/logz"
)

type AppContext struct {
	components []IComponent
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

// InitComponent 添加组件, 并触发初始化
func (x *AppContext) InitComponent(components ...IComponent) error {
	x.components = make([]IComponent, 0, len(components))
	x.components = append(x.components, components...)

	err := x.call("init", func(c IComponent) error {
		return c.Init()
	})
	if err != nil {
		logz.Fatal().Msgf("init component error %+v", err)
	}
	return err
}

func (x *AppContext) Start() error {
	err := x.call("start", func(c IComponent) error {
		return c.Start()
	})
	if err != nil {
		logz.Fatal().Msgf("start component error %+v", err)
	}
	return err
}

func (x *AppContext) Showdown(wait int) {
	_ = x.callReverse("showdown", func(c IComponent) error {
		// 执行关闭回调函数
		err := c.Showdown()
		if err != nil {
			logz.Error().Msgf("%v shutdown err: %v", c.Name(), err)
		}
		return nil
	})

	// 等待, 避免有未处理完的异步事务
	for i := wait; i >= 0; i-- {
		logz.Info().Msgf("Close Wait: %v", i)
		time.Sleep(time.Second)
	}
	logz.Info().Msg("Finish Close Wait")
}

// MainRun 主程序使用的组合调用;
// 结合了 InitComponent, Start, WaitExitSig, Showdown
func (x *AppContext) MainRun(wait int, components ...IComponent) {
	defer x.Showdown(wait)

	// init
	if x.InitComponent(components...) != nil {
		return
	}

	// start
	if x.Start() != nil {
		return
	}

	// 等待结束信号(service 处理收到信号后, 不再会有新的请求打过来)
	WaitExitSig()
}

func (x *AppContext) call(callName string, f func(c IComponent) error) error {
	start := time.Now()
	for _, component := range x.components {
		st := time.Now()
		err := f(component)
		total := time.Since(st)
		if err != nil {
			logz.Error().Msgf("%v %v error:[%v] time:%vms", callName, component.Name(), err, total.Milliseconds())
			return err
		}
		logz.Info().Msgf("%v %v success; time:%vms", callName, component.Name(), total.Milliseconds())
	}
	logz.Info().Msgf("finish %v all components time:%vms", callName, time.Since(start).Milliseconds())
	return nil
}

func (x *AppContext) callReverse(callName string, f func(c IComponent) error) error {
	start := time.Now()
	for i := len(x.components) - 1; i >= 0; i-- {
		component := x.components[i]
		st := time.Now()
		err := f(component)
		total := time.Since(st)
		if err != nil {
			logz.Error().Msgf("%v %v error:[%v] time:%vms", callName, component.Name(), err, total.Milliseconds())
			return err
		}
		logz.Info().Msgf("%v %v success; time:%vms", callName, component.Name(), total.Milliseconds())
	}
	logz.Info().Msgf("finish %v all components time:%vms", callName, time.Since(start).Milliseconds())
	return nil
}
