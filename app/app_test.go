package app

import (
	"fmt"
	"github.com/geeklingithub/best-go/server/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	//业务http
	businessHttp := http.New(
		http.Address("127.0.0.1:8088"),
		http.ShutdownFunc(func() {
			fmt.Println("优雅关服的一些钩子操作,将缓存中的数据回写到数据库中")
			time.Sleep(time.Minute)
		}),
	)
	//添加路由
	businessHttp.AddRouter(http.BusinessRouter())

	//管理http
	adminHttp := http.New(
		http.Address("127.0.0.1:8089"),
		http.ShutdownFunc(func() {
			fmt.Println("优雅关服的一些钩子操作,将缓存中的数据回写到数据库中")
			time.Sleep(time.Minute)
		}),
	)
	//添加路由
	adminHttp.AddRouter(http.AdminRouter())

	//新建应用
	app := New(
		Name("应用名称"),
		Version("v1.0"),
		Servers(businessHttp, adminHttp),
		//Servers(businessHttp),
	)

	//启动应用
	app.Start()
	time.AfterFunc(20*time.Second, func() {
		app.Stop()
	})
	time.Sleep(time.Hour)
}
