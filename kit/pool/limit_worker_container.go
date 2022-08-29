package pool

import "time"

type workerLoopQueue struct {
	running  []*Worker
	expired  []*Worker
	head     int  //头
	tail     int  //尾
	capacity int  //容量
	isFull   bool //是否已满
}

// newWorkerLoopQueue 新建循环队列
func newWorkerLoopQueue(capacity int) *workerLoopQueue {
	return &workerLoopQueue{
		running:  make([]*Worker, capacity),
		capacity: capacity,
	}
}

// size 队列当前数量
func (wq *workerLoopQueue) size() int {
	if wq.capacity == 0 {
		return 0
	}

	if wq.head == wq.tail {
		if wq.isFull {
			return wq.capacity
		}
		return 0
	}

	if wq.tail > wq.head {
		return wq.tail - wq.head
	}

	return wq.capacity - wq.head + wq.tail
}

func (wq *workerLoopQueue) isEmpty() bool {
	return wq.head == wq.tail && !wq.isFull
}

func (wq *workerLoopQueue) append(worker *Worker) error {

	if wq.capacity == 0 {
		return errWorkerContainerIsClosed
	}

	if wq.isFull {
		return errWorkerContainerIsFull
	}

	wq.running[wq.tail] = worker
	wq.tail++

	if wq.tail == wq.capacity {
		wq.tail = 0
	}
	if wq.tail == wq.head {
		wq.isFull = true
	}

	return nil
}

func (wq *workerLoopQueue) pop() *Worker {
	if wq.isEmpty() {
		return nil
	}

	w := wq.running[wq.head]
	wq.running[wq.head] = nil
	wq.head++
	if wq.head == wq.capacity {
		wq.head = 0
	}
	wq.isFull = false

	return w
}

func (wq *workerLoopQueue) recycleExpiredWorker(duration time.Duration) []*Worker {
	expiryTime := time.Now().Add(-duration)
	index := wq.binarySearch(expiryTime)
	if index == -1 {
		return nil
	}
	wq.expired = wq.expired[:0]

	if wq.head <= index {
		wq.expired = append(wq.expired, wq.running[wq.head:index+1]...)
		for i := wq.head; i < index+1; i++ {
			wq.running[i] = nil
		}
	} else {
		wq.expired = append(wq.expired, wq.running[0:index+1]...)
		wq.expired = append(wq.expired, wq.running[wq.head:]...)
		for i := 0; i < index+1; i++ {
			wq.running[i] = nil
		}
		for i := wq.head; i < wq.capacity; i++ {
			wq.running[i] = nil
		}
	}
	head := (index + 1) % wq.capacity
	wq.head = head
	if len(wq.expired) > 0 {
		wq.isFull = false
	}

	return wq.expired
}

func (wq *workerLoopQueue) binarySearch(expiryTime time.Time) int {
	var mid, size, basel, tempMid int
	size = len(wq.running)

	// if no need to remove work, return -1
	if wq.isEmpty() || expiryTime.Before(wq.running[wq.head].recycleTime) {
		return -1
	}

	// example
	// size = 8, head = 7, tail = 4
	// [ 2, 3, 4, 5, nil, nil, nil,  1]
	//   0  1  2  3    4   5     6   7
	//              tail          head
	//
	//   1  2  3  4  nil nil   nil   0
	//            r                  l

	r := (wq.tail - 1 - wq.head + size) % size
	basel = wq.head
	l := 0
	for l <= r {
		mid = l + ((r - l) >> 1)
		tempMid = (mid + basel + size) % size
		if expiryTime.Before(wq.running[tempMid].recycleTime) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return (r + basel + size) % size
}

func (wq *workerLoopQueue) reset() {
	if wq.isEmpty() {
		return
	}

	for w := wq.pop(); w != nil; {
		w.taskFunc <- nil
	}
	wq.running = wq.running[:0]
	wq.capacity = 0
	wq.head = 0
	wq.tail = 0

}
