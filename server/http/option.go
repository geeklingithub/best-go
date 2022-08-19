package http

type Option struct {
	address string

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
