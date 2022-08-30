package pool

import "context"

type OrderPool struct {
	workers      []*worker
	maxWaitTasks int32
}

//NewOrderPool 新建有序任务池
func NewOrderPool(workerSize int, maxWaitTasks int32) *OrderPool {

	workers := make([]*worker, workerSize, workerSize)
	pool := &OrderPool{
		workers:      workers,
		maxWaitTasks: maxWaitTasks,
	}
	for i := 0; i < len(workers); i++ {
		workers[i] = &worker{
			pool:            pool,
			taskFuncChannel: make(chan *WorkerTask, maxWaitTasks),
			currencyNum:     maxWaitTasks,
		}
	}
	return pool
}

//Shutdown 优雅关闭
func (pool *OrderPool) Shutdown() {

}

//ShutdownNow 立即关闭
func (pool *OrderPool) ShutdownNow() {

}

//WaitingShutdown 等待关闭
func (pool *OrderPool) WaitingShutdown() {

}

//SubmitTask 提交任务
func (pool *OrderPool) SubmitTask(workerKey int, taskFunc func(context.Context)) (context.Context, error) {
	workerLen := len(pool.workers)
	return pool.workers[workerKey%workerLen].submitTask(taskFunc)
}

//SubmitScheduleTask 提交调度任务
func (pool *OrderPool) SubmitScheduleTask(workerKey int, taskFunc func(context.Context)) (context.Context, error) {
	workerLen := len(pool.workers)
	return pool.workers[workerKey%workerLen].submitTask(taskFunc)
}
