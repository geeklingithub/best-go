package app

import (
	"log"
	"os"
	"time"
)

// AppOption 配置项
type AppOption struct {
	name            string        //应用名称
	version         string        //应用版本
	logger          log.Logger    //应用日志
	servers         []Server      //应用服务列表
	shutdownTimeOut time.Duration //应用优雅关服超时时间
	closeSignals    []os.Signal   //应用关服信号
}

type OptFunc func(option *AppOption)

func Name(name string) OptFunc {
	return func(option *AppOption) {
		option.name = name
	}
}

func Version(version string) OptFunc {
	return func(option *AppOption) {
		option.version = version
	}
}

func Servers(servers ...Server) OptFunc {
	return func(option *AppOption) {
		option.servers = servers
	}
}
