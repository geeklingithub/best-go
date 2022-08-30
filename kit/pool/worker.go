package pool

import (
	"errors"
	"sync/atomic"
)

const (
	CLOSED = 0
	OPENED = 1
)

var (
	ErrOverloadTaskSize = errors.New("OverloadTaskSize")
)

type worker struct {
	pool            *OrderPool       // 工人所属池子
	taskFuncChannel chan *WorkerTask // 任务
	currencyNum     int32            //并发数
	runTaskNum      int32            // 状态
}

type WorkerTask struct {
	taskFunc func()
}

//alterRunTaskNum 修改令牌数
func (worker *worker) alterRunTaskNum(changeNum int) int32 {
	return atomic.AddInt32(&worker.runTaskNum, int32(changeNum))
}

//alterRunTaskNum 修改令牌数
func (worker *worker) getRunTaskNum() int32 {
	return atomic.LoadInt32(&worker.runTaskNum)
}

//submitTask 提交任务
func (worker *worker) submitTask(taskFunc func()) error {

	if worker.getRunTaskNum() >= worker.currencyNum {
		return ErrOverloadTaskSize
	}

	runTaskNum := worker.alterRunTaskNum(1)
	worker.taskFuncChannel <- &WorkerTask{
		taskFunc: taskFunc,
	}

	if runTaskNum == OPENED {
		go func() {
			defer func() {

			}()

			for workerTask := range worker.taskFuncChannel {
				if workerTask == nil {
					return
				}
				workerTask.taskFunc()
				worker.alterRunTaskNum(-1)
				if worker.getRunTaskNum() == CLOSED {
					return
				}
			}
		}()
	}
	return nil
}
