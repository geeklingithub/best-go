package http

import "net/http"

type Option struct {
	address   string
	routerMap map[string]func(writer http.ResponseWriter, request *http.Request)
	shutdown  func()
}

type OptFunc func(*Option)

func Address(address string) OptFunc {
	return func(option *Option) {
		option.address = address
	}
}

func RouterMap(routerMap map[string]func(writer http.ResponseWriter, request *http.Request)) OptFunc {
	return func(option *Option) {
		option.routerMap = routerMap
	}
}

func Shutdown(shutdownFunc func()) OptFunc {
	return func(option *Option) {
		option.shutdown = shutdownFunc
	}
}
