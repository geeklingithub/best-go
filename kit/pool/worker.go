package pool

import (
	"time"
)

type Worker struct {
	pool        *NoOrderCoroutinePool // 工人所属池子
	taskFunc    chan func()           // 任务
	recycleTime time.Time             //工人回收时间
}

func (worker *Worker) exec() {
	worker.pool.AlterRunning(1)
	go func() {
		defer func() {
			worker.pool.AlterRunning(-1)
			worker.pool.workerCache.Put(worker)
		}()

		for f := range worker.taskFunc {
			if f == nil {
				return
			}
			f()
			if ok := worker.pool.recycleWorker(worker); !ok {
				return
			}
		}
	}()
}
