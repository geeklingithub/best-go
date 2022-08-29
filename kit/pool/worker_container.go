package pool

import (
	"errors"
	"time"
)

var (
	// errWorkerContainerIsFull 队列满
	errWorkerContainerIsFull = errors.New("the queue is full")

	// errWorkerContainerIsClosed 协程池已关闭
	errWorkerContainerIsClosed = errors.New("the queue length is zero")
)

type WorkerContainer interface {
	size() int
	isEmpty() bool
	append(worker *Worker) error
	pop() *Worker
	recycleExpiredWorker(duration time.Duration) []*Worker
	reset()
}

type workerContainerType int

const (
	limitType workerContainerType = 1 << iota
	unLimitType
)

func newWorkerArray(workerContainerType workerContainerType, size int) WorkerContainer {
	switch workerContainerType {
	//case unLimitType:
	//	return newWorkerStack(size)
	case limitType:
		return newWorkerLoopQueue(size)
	default:
		//return newWorkerStack(size)
		return nil
	}
}
