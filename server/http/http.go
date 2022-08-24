package best_http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
	*Option
	ctx      context.Context
	cancel   context.CancelFunc
	shutdown bool
	router   *Router
}

// New 服务初始化
func New(opts ...OptFunc) *Server {
	//初始化配置项
	//默认配置
	o := &Option{}

	//自定义配置
	for _, opt := range opts {
		opt(o)
	}

	//返回应用实例对象
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		server:   &http.Server{},
		Option:   o,
		cancel:   cancel,
		ctx:      ctx,
		shutdown: false,
		router:   &Router{methodMap: map[string]*RouterInfo{}},
	}
}

// Start 服务启动
func (server *Server) Start(context.Context) error {

	for routerPath, _ := range server.router.methodMap {
		http.HandleFunc(routerPath, server.HandleFunc(routerPath))
	}
	fmt.Println("http 启动 ", server.Option.address)
	server.server.Addr = server.Option.address
	return server.server.ListenAndServe()
}

func (server *Server) Stop(ctx context.Context) error {
	fmt.Println("http 开始优雅关闭 ", server.Option.address)
	server.shutdown = true
	server.shutdownFunc()
	err := server.server.Shutdown(ctx)
	return err
}

func (server *Server) HandleFunc(key string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		if server.shutdown {
			writer.Write([]byte("shutdown中,拒绝请求"))
			return
		}
		len := request.ContentLength
		body := make([]byte, len)
		// 将请求体中内容读到body中
		request.Body.Read(body)
		routerInfo := server.router.methodMap[key]
		json.Unmarshal(body, routerInfo.reqBody)

		ctx := NewContext{
			writer:  writer,
			Request: request,
		}
		routerInfo.handleFunc(routerInfo.reqBody, ctx)
	}
}

func (server *Server) AddRouter(f func(router *Router) error) error {
	return f(server.router)
}
