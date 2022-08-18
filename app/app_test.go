package app

import (
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	app := Init(
		Name("应用名称"),
		Version("v1.0"),
	)
	app.Start()
	//time.AfterFunc(time.Second, func() {
	//	app.Stop()
	//})
	time.Sleep(time.Minute)
}
