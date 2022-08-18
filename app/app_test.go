package app

import (
	"github.com/geeklingithub/best-go/server/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	businessHttp := http.New(
		http.Address("127.0.0.1:8080"),
		http.RouterMap(http.BusinessRouter()),
	)
	adminHttp := http.New(
		http.Address("127.0.0.1:8081"),
		http.RouterMap(http.AdminRouter()),
	)
	app := New(
		Name("应用名称"),
		Version("v1.0"),
		Servers(businessHttp, adminHttp),
	)
	app.Start()
	time.AfterFunc(time.Second, func() {
		app.Stop()
	})
	time.Sleep(time.Minute)
}
