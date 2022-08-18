package http

import "net/http"

func BusinessRouter() map[string]func(writer http.ResponseWriter, request *http.Request) {
	return map[string]func(writer http.ResponseWriter, request *http.Request){
		"/login": func(writer http.ResponseWriter, request *http.Request) {

		},
		"/createPlayer": func(writer http.ResponseWriter, request *http.Request)) {

	},
		"/selectGateway": func(writer http.ResponseWriter, request *http.Request) {

	},
	}
}

func AdminRouter() map[string]func(writer http.ResponseWriter, request *http.Request) {
	return map[string]func(writer http.ResponseWriter, request *http.Request){
		"/admin": func(writer http.ResponseWriter, request *http.Request) {

		},
	}
}