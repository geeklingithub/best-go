package pool

import (
	"errors"
	"time"
)

var (
	// errWorkerContainerIsFull 满
	errWorkerContainerIsFull = errors.New("the capacity is full")

	// errPoolIsClosed 协程池已关闭
	errPoolIsClosed = errors.New("the queue length is zero")
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
)

func NewWorkerContainer(size int) WorkerContainer {
	return newWorkerLoopQueue(size)
}
