package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// NoOrderCoroutinePool 无序协程池
type NoOrderCoroutinePool struct {
	options         *Option         //配置项
	capacity        int32           //池子容量
	state           int32           //池子状态
	workerContainer WorkerContainer //worker容器
	workerCache     sync.Pool       //缓存池，加速worker的获取
	running         int32           //正在运行的goroutine数量

}

// 协程池状态
const (
	closed = iota
	open
)

// popIdleWorker 获取一个空闲的worker
func (pool *NoOrderCoroutinePool) popIdleWorker() (worker *Worker) {
	worker = pool.workerContainer.pop()
	return
}

// recycleWorker 回收worker
func (pool *NoOrderCoroutinePool) recycleWorker(worker *Worker) bool {

	worker.recycleTime = time.Now()
	pool.workerContainer.append(worker)
	return true
}

// cleanExpiredWorkerDuration 定期清理过期的worker
func (pool *NoOrderCoroutinePool) cleanExpiredWorkerDuration(ctx context.Context) {
	durationCheck := time.NewTicker(pool.options.WorkerExpiredDuration)
	defer func() {
		durationCheck.Stop()
	}()
	for {
		select {
		//定时唤醒
		case <-durationCheck.C:
			//释放
		case <-ctx.Done():
			return
		}

		//获取过期的worker
		expiredWorkers := pool.workerContainer.recycleExpiredWorker(pool.options.WorkerExpiredDuration)
		for i := range expiredWorkers {
			expiredWorkers[i].taskFunc <- nil
			expiredWorkers[i] = nil
		}

	}
}

//AlterRunning 修改goroutine数量
func (pool *NoOrderCoroutinePool) AlterRunning(changeNum int) {
	atomic.AddInt32(&pool.running, int32(changeNum))
}

// Close 关闭协程池并释放worker
func (pool *NoOrderCoroutinePool) Close() {
	if !atomic.CompareAndSwapInt32(&pool.state, open, closed) {
		return
	}
	pool.workerContainer.reset()
}
