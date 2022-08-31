package pool

import (
	"context"
	"errors"
	"sync/atomic"
)

const (
	GoClosed      = 0 //消费任务goroutine处于关闭状态
	GoOpened      = 1 //消费任务goroutine处于开启状态
	WorkerClosing = 2 //worker处于关闭中状态,不在提供服务
	WorkerClosed  = 3 //worker处于完全关闭状态
)

var (
	ErrOverloadTaskSize = errors.New("OverloadTaskSize")
	ErrWorkerClose      = errors.New("WorkerClose")
)

type worker struct {
	pool            *OrderPool       // 工人所属池子
	taskFuncChannel chan *FutureTask // 任务
	currencyNum     int32            //最大并发数
	runTaskNum      int32            // 当前任务数
	state           int32            // 消费任务的goroutine状态
	ctx             context.Context  //worker上下文
	cancel          context.CancelFunc
	closedSignal    chan bool //退出worker channel通信
}

type WorkerTaskResult struct {
	result any
}

//alterRunTaskNum 修改任务数
func (worker *worker) alterRunTaskNum(changeNum int) int32 {
	return atomic.AddInt32(&worker.runTaskNum, int32(changeNum))
}

//getRunTaskNum 当前任务数
func (worker *worker) getRunTaskNum() int32 {
	return atomic.LoadInt32(&worker.runTaskNum)
}

//submitTask 提交任务
func (worker *worker) submitTask(taskFunc func() any) (*FutureTask, error) {

	if atomic.LoadInt32(&worker.state) == WorkerClosed {
		return nil, ErrWorkerClose
	}

	if atomic.LoadInt32(&worker.state) == WorkerClosing {
		return nil, ErrWorkerClose
	}

	if worker.getRunTaskNum() >= worker.currencyNum {
		return nil, ErrOverloadTaskSize
	}

	worker.alterRunTaskNum(1)
	futureTaskCtx, futureTaskCancel := context.WithCancel(worker.ctx)
	futureTask := &FutureTask{
		taskFunc: taskFunc,
		ctx:      futureTaskCtx,
		cancel:   futureTaskCancel,
	}
	worker.taskFuncChannel <- futureTask

	// state 保证统一时刻只有一个goroutine消费任务
	if atomic.CompareAndSwapInt32(&worker.state, GoClosed, GoOpened) {
		worker.run()
	}
	return futureTask, nil
}

func (worker *worker) run() {
	go func() {
		defer func() {
			atomic.CompareAndSwapInt32(&worker.state, GoOpened, GoClosed)
		}()

		//消费任务
		for futureTask := range worker.taskFuncChannel {

			if worker.isClosed() {
				return
			}

			if futureTask == nil {
				return
			}
			futureTask.taskResult = futureTask.taskFunc()
			futureTask.cancel()
			taskNum := worker.alterRunTaskNum(-1)
			if taskNum == 0 {
				return
			}
		}
	}()
}

//shutdown 优雅关闭
func (worker *worker) shutdown() {
	//worker打上关闭标志,不在接受新任务
	atomic.StoreInt32(&worker.state, WorkerClosing)
	//等待退出
	worker.waitingShutdown()
}

//waitingShutdown 等待关闭
func (worker *worker) waitingShutdown() {
	go func() {
		ctx, cancel := context.WithTimeout(worker.ctx, worker.pool.shutdownTimeout)
		defer cancel()
		select {
		case <-ctx.Done():
			//通知worker退出
			worker.shutdownNow()
		}
	}()
}

//waitingShutdown 等待关闭
func (worker *worker) shutdownNow() {
	//通知worker退出
	worker.closedSignal <- true
	worker.closedSignal = nil
}

//isClosed 是否关闭
func (worker *worker) isClosed() bool {

	if atomic.LoadInt32(&worker.state) == WorkerClosed {
		return true
	}

	select {
	//因取消或超时,退出worker执行
	case <-worker.closedSignal:
		//worker打上关闭标志,不在接受新任务
		atomic.StoreInt32(&worker.state, WorkerClosed)
		//todo
		// 记录未执行任务现场
		// 处理死循环
		return true
	default:
		return false
	}
}
