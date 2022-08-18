package http

import "net/http"

func BusinessRouter() map[string]func(writer http.ResponseWriter, request *http.Request) {
	return map[string]func(writer http.ResponseWriter, request *http.Request){
		"/login": func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("/login"))
		},
		"/createPlayer": func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("/createPlayer"))
		},
		"/selectGateway": func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("/selectGateway"))
		},
	}
}

func AdminRouter() map[string]func(writer http.ResponseWriter, request *http.Request) {
	return map[string]func(writer http.ResponseWriter, request *http.Request){
		"/admin": func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("/admin"))
		},
	}
}
