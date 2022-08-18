package http

import "net/http"

type Option struct {
	address   string
	routerMap map[string]func(writer http.ResponseWriter, request *http.Request)
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
