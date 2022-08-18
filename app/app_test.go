package app

import (
	"github.com/geeklingithub/best-go/server/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	businessHttp := http.New(
		http.Address("127.0.0.1:8088"),
		http.RouterMap(http.BusinessRouter()),
	)
	adminHttp := http.New(
		http.Address("127.0.0.1:8089"),
		http.RouterMap(http.AdminRouter()),
	)
	app := New(
		Name("应用名称"),
		Version("v1.0"),
		Servers(businessHttp, adminHttp),
	)
	app.Start()
	time.AfterFunc(time.Minute, func() {
		app.Stop()
	})
	time.Sleep(time.Hour)
}
