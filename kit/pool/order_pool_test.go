package pool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewOrderPool(t *testing.T) {
	pool := NewOrderPool(1, 100)

	for i := 0; i < 10000; i++ {
		index := i

		ctx, err := pool.SubmitTask(index, func(ctx context.Context) {
			fmt.Println("count", index, index%100)
		})
		if err != nil {
			return
		}

		fmt.Println(ctx, err)
	}

	time.Sleep(time.Minute)
}
