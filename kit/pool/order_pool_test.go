package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewOrderPool(t *testing.T) {
	pool := NewOrderPool(100, 100, 5*time.Minute)

	for i := 0; i < 10000; i++ {
		index := i

		futureTask, err := pool.SubmitTask(index, func() any {
			fmt.Println("count", index, index%100)
			return index
		})
		if err != nil {
			return
		}
		go func() {
			select {
			case <-futureTask.ctx.Done():
				fmt.Println("futureTask", futureTask.taskResult)
			}
		}()

	}

	time.Sleep(time.Minute)
}
