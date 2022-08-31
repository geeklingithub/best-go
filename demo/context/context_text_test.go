package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextCancel(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("取消")
		}
	}()

	cancel()
	cancel()
	cancel()
	cancel()
	time.Sleep(time.Minute)
}
