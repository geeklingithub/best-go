package pool

import (
	"context"
	"time"
)

type OrderPool struct {
	workers         []*worker
	maxWaitTasks    int32
	shutdownTimeout time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
}

//NewOrderPool 新建有序任务池
func NewOrderPool(workerSize int, maxWaitTasks int32, shutdownTimeout time.Duration) *OrderPool {

	workers := make([]*worker, workerSize, workerSize)
	ctx, cancel := context.WithCancel(context.Background())
	pool := &OrderPool{
		workers:         workers,
		maxWaitTasks:    maxWaitTasks,
		ctx:             ctx,
		cancel:          cancel,
		shutdownTimeout: shutdownTimeout,
	}

	for i := 0; i < len(workers); i++ {
		workerCtx, workerCancel := context.WithCancel(ctx)
		workers[i] = &worker{
			pool:            pool,
			taskFuncChannel: make(chan *WorkerTask, maxWaitTasks),
			currencyNum:     maxWaitTasks,
			ctx:             workerCtx,
			cancel:          workerCancel,
		}
	}
	return pool
}

//Shutdown 优雅关闭
func (pool *OrderPool) Shutdown() {
	for _, worker := range pool.workers {
		worker.shutdown()
	}
}

//ShutdownNow 强制关闭
func (pool *OrderPool) ShutdownNow() {
	for _, worker := range pool.workers {
		worker.shutdownNow()
	}
}

//SubmitTask 提交任务
func (pool *OrderPool) SubmitTask(workerKey int, taskFunc func(context.Context)) (context.Context, error) {
	workerLen := len(pool.workers)
	return pool.workers[workerKey%workerLen].submitTask(taskFunc)
}
