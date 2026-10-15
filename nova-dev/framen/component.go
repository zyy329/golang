package framen

type IComponent interface {
	Name() string
	Init() error
	Start() error
	Showdown() error
}

type SimpleComponent struct {
	initS []func()
}

//goland:noinspection GoUnusedExportedFunction
func NewSimpleComponent(initS ...func()) *SimpleComponent {
	return &SimpleComponent{initS: initS}
}

func (c *SimpleComponent) Name() string {
	return ""
}
func (c *SimpleComponent) Init() error {
	for _, f := range c.initS {
		f()
	}
	return nil
}
func (c *SimpleComponent) Start() error {
	return nil
}
func (c *SimpleComponent) Showdown() error {
	return nil
}
