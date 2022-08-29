package pool

import (
	"context"
	"errors"
	"github.com/geeklingithub/best-go/kit/queue"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (

	// ErrInvalidWorkerExpiredDuration 清理worker时的间隔参数未负数
	ErrInvalidWorkerExpiredDuration = errors.New(" worker expired duration for pool must be positive ")

	// ErrPoolClosed 池子关闭之后，提交任务
	ErrPoolClosed = errors.New("this pool has been closed")
)

type taskFunc chan func()

// NoOrderCoroutinePool 有序协程池
type NoOrderCoroutinePool struct {
	options         *Options                           //配置项
	capacity        int32                              //池子容量
	state           int32                              //池子状态
	workerContainer WorkerContainer                    //worker容器
	workerCache     sync.Pool                          //缓存池，加速worker的获取
	running         int32                              //正在运行的goroutine数量
	taskFuncQueue   queue.ArrayBlockingQueue[taskFunc] //阻塞的任务队列,这里用队列是为了公平
	stopFunc        context.CancelFunc
}

type blockTask struct {
	waitTaskFuncQueue []*taskFunc //任务
	position          int         //当前位置
	cond              *sync.Cond  //等待获取一个空闲的worker
}

// 协程池状态
const (
	closed = iota
	open
)

var (
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	workerChanCap = func() int {
		// Use blocking channel if GOMAXPROCS=1.
		// This switches context from sender to receiver immediately,
		// which results in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the sender might be dragged down if the receiver is CPU-bound.
		return 1
	}()
)

// NewNoOrderCoroutinePool 创捷协程池实例
func NewNoOrderCoroutinePool(size int, options ...OptionFunc) (*NoOrderCoroutinePool, error) {
	opts, err := initOptions(options...)

	if err != nil {
		return nil, err
	}

	if size <= 0 {
		size = -1
	}

	pool := &NoOrderCoroutinePool{
		capacity: int32(size),
		options:  opts,
	}
	pool.workerCache.New = func() interface{} {
		return &Worker{
			pool:     pool,
			taskFunc: make(chan func(), workerChanCap),
		}
	}

	pool.workerContainer = NewWorkerContainer(size)

	// 启动一个goroutine去定时清理过期的worker
	var ctx context.Context
	ctx, pool.stopFunc = context.WithCancel(context.Background())
	go pool.cleanExpiredWorkerDuration(ctx)

	return pool, nil
}

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

func (pool *NoOrderCoroutinePool) Submit(taskFunc func()) error {
	if atomic.LoadInt32(&pool.state) == closed {
		return ErrPoolClosed
	}
	worker := pool.popIdleWorker()
	worker.taskFunc <- taskFunc
	return nil
}
