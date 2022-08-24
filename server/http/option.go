package best_http

type Option struct {
	address      string
	filterChain  Filter
	shutdownFunc func()
}

type OptFunc func(*Option)

func Address(address string) OptFunc {
	return func(option *Option) {
		option.address = address
	}
}

func ShutdownFunc(shutdownFunc func()) OptFunc {
	return func(option *Option) {
		option.shutdownFunc = shutdownFunc
	}
}

func FilterChain(fs ...FilterHandle) OptFunc {
	return func(option *Option) {
		option.filterChain = option.filterChain.AddFilter(fs...)
	}
}
