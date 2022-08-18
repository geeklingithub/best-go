package http

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
	*Option
	ctx    context.Context
	cancel context.CancelFunc
}

// Init 服务初始化
func New(opts ...OptFunc) Server {
	//初始化配置项
	//默认配置
	o := &Option{}

	//自定义配置
	for _, opt := range opts {
		opt(o)
	}

	//返回应用实例对象
	ctx, cancel := context.WithCancel(context.Background())
	return Server{
		server: &http.Server{},
		Option: o,
		cancel: cancel,
		ctx:    ctx,
	}
}

// Start 服务启动
func (server Server) Start(context.Context) error {

	for routerPath, handleFunc := range server.Option.routerMap {
		http.HandleFunc(routerPath, handleFunc)
	}
	fmt.Println("http 启动 ", server.Option.address)
	server.server.Addr = server.Option.address
	return server.server.ListenAndServe()
}

// Stop 服务关闭
func (server Server) Stop(ctx context.Context) error {
	fmt.Println("http 关闭 ", server.Option.address)
	return server.server.Shutdown(ctx)
}
