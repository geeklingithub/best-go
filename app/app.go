package app

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

// App 应用
type App struct {
	*Option                    //应用配置项
	ctx     context.Context    //应用上下文
	cancel  context.CancelFunc //上下文取消信号
}

// Server 服务接口
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// New 应用创建
func New(opts ...OptFunc) *App {
	//初始化配置项

	//默认配置
	o := &Option{
		closeSignals: []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT},
	}

	//自定义配置
	for _, opt := range opts {
		opt(o)
	}

	//返回应用实例对象
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		Option: o,
		cancel: cancel,
		ctx:    ctx,
	}
}

// Start 应用启动
func (app *App) Start() {

	wg := &sync.WaitGroup{}
	for _, server := range app.servers {
		wg.Add(1)

		//服务关闭
		go func() {
			//应用关闭时,关闭服务
			<-app.ctx.Done()
			ctx, cancel := context.WithTimeout(app.ctx, app.stopTimeOut)
			defer cancel()
			server.Stop(ctx)
		}()

		//服务启动
		server := server
		go func() {
			wg.Done()
			server.Start(app.ctx)
		}()
	}

	wg.Wait()

	//信号通知关服
	c := make(chan os.Signal, 1)
	signal.Notify(c, app.closeSignals...)
	go func() {
		select {
		//非信号退出时,及时回收goroutine
		case <-app.ctx.Done():
		//信号退出时,优雅关闭应用
		case <-c:

			//再次收到信号，强制退出
			force := make(chan os.Signal, 1)
			signal.Notify(c, app.closeSignals...)
			go func() {
				<-force
				runtime.Goexit()
			}()
			//优雅关闭应用
			app.Stop()
		}
	}()
}

// Stop 应用关闭
func (app *App) Stop() {

	app.cancel()
}
