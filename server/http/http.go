package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
	*Option
	ctx       context.Context
	cancel    context.CancelFunc
	shutdown  bool
	routerMap map[string]func(writer http.ResponseWriter, request *http.Request)
}

// Init 服务初始化
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
		server:    &http.Server{},
		Option:    o,
		cancel:    cancel,
		ctx:       ctx,
		shutdown:  false,
		routerMap: map[string]func(writer http.ResponseWriter, request *http.Request){},
	}
}

// Start 服务启动
func (server *Server) Start(context.Context) error {

	for routerPath, handleFunc := range server.routerMap {
		http.HandleFunc(routerPath, handleFunc)
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

func (server *Server) AddRouter(routerMap map[string]func() any) {
	for key, value := range routerMap {
		server.routerMap[key] = server.HandleFunc(value)
	}
}

func (server *Server) HandleFunc(f func() any) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if server.shutdown {
			writer.Write([]byte("shutdown中,拒绝请求"))
			return
		}
		resp := f()
		v, _ := json.Marshal(resp)
		writer.Write(v)
	}
}
