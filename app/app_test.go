package app

import (
	"fmt"
	"github.com/geeklingithub/best-go/business/login"
	best_http "github.com/geeklingithub/best-go/server/http"

	"testing"
	"time"
)

func TestApp(t *testing.T) {
	//业务http
	businessHttp := best_http.New(
		best_http.Address("127.0.0.1:8088"),
		best_http.ShutdownFunc(func() {
			fmt.Println("优雅关服的一些钩子操作,将缓存中的数据回写到数据库中")
			time.Sleep(time.Minute)
			fmt.Println("优雅关服完毕")
		}),
		//best_http.FilterChain(best_http.MetricFilterHandle),
	)
	businessHttp.AddRouter(login.RegisterRouter)

	//管理http
	//adminHttp := best_http.New(
	//	best_http.Address("127.0.0.1:8089"),
	//	best_http.ShutdownFunc(func() {
	//		fmt.Println("优雅关服的一些钩子操作,将缓存中的数据回写到数据库中")
	//		time.Sleep(time.Minute)
	//		fmt.Println("优雅关服完毕")
	//	}),
	//)

	//新建应用
	app := New(
		Name("应用名称"),
		Version("v1.0"),
		//Servers(businessHttp, adminHttp),
		Servers(businessHttp),
	)

	//启动应用
	app.Start()
	//time.AfterFunc(5*time.Second, func() {
	//	app.Stop()
	//
	//})
	time.Sleep(3 * time.Minute)
}
