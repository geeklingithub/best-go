package app

import (
	"fmt"
	"github.com/geeklingithub/best-go/server/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	businessHttp := http.New(
		http.Address("127.0.0.1:8088"),
		http.RouterMap(http.BusinessRouter()),
		http.Shutdown(func() {
			fmt.Println("优雅关服的一些钩子操作,将缓存中的数据回写到数据库中")
		}),
	)
	adminHttp := http.New(
		http.Address("127.0.0.1:8089"),
		http.RouterMap(http.AdminRouter()),
		http.Shutdown(func() {
			fmt.Println("优雅关服的一些钩子操作,将缓存中的数据回写到数据库中")
		}),
	)
	app := New(
		Name("应用名称"),
		Version("v1.0"),
		Servers(businessHttp, adminHttp),
	)
	app.Start()
	time.AfterFunc(10*time.Second, func() {
		app.Stop()
	})
	time.Sleep(time.Hour)
}
