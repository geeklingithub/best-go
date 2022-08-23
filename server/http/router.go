package best_http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Router struct {
	methodMap map[string]*RouterInfo
}

type RouterInfo struct {
	reqBody    any
	handleFunc func(any, NewContext)
}

type NewContext struct {
	writer  http.ResponseWriter
	Request *http.Request
}

func (ctx NewContext) SendResp(resp any) {
	respJson, _ := json.Marshal(resp)
	ctx.writer.Write(respJson)
}

func (router Router) AddRouter(path string, fun any, reqBody any) {
	key := path
	_, ok := router.methodMap[key]
	if ok {
		fmt.Println("repeat router")
		return
	}
	router.methodMap[key] = &RouterInfo{
		reqBody: reqBody,
		handleFunc: func(reqBody any, ctx NewContext) {
			args := make([]reflect.Value, 0, 2)
			args = append(args, reflect.ValueOf(reqBody))
			args = append(args, reflect.ValueOf(ctx))
			reflect.ValueOf(fun).Call(args)
		},
	}
}
