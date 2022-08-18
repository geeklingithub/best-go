package app

import (
	"log"
	"os"
	"time"
)

// Option 配置项
type Option struct {
	name         string        //应用名称
	version      string        //应用版本
	logger       log.Logger    //应用日志
	servers      []Server      //应用服务列表
	stopTimeOut  time.Duration //应用优雅关服超时时间
	closeSignals []os.Signal   //应用关服信号
}

type OptFunc func(option *Option)

func Name(name string) OptFunc {
	return func(option *Option) {
		option.name = name
	}
}

func Version(version string) OptFunc {
	return func(option *Option) {
		option.version = version
	}
}
