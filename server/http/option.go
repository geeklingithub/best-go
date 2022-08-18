package http

type Option struct {
	port      int
	routerMap map[string]any
}

type OptFunc func(*Option)

func (o *Option) Port(port int) OptFunc {
	return func(option *Option) {
		option.port = port
	}
}
