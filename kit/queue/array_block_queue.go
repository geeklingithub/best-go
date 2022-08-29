package queue

import "sync"

type ArrayBlockingQueue[T any] struct {
	items []T

	head int

	tail int

	count    int
	notEmpty sync.Cond
	notFull  sync.Cond
}

func (queue *ArrayBlockingQueue[T]) NewArrayBlockingQueue(size int) *ArrayBlockingQueue[T] {
	return &ArrayBlockingQueue[T]{
		items:    make([]T, size),
		head:     0,
		tail:     0,
		count:    0,
		notEmpty: sync.Cond{},
		notFull:  sync.Cond{},
	}
}

func (queue *ArrayBlockingQueue[T]) Enqueue(e T) {

	for queue.count == len(queue.items) {
		queue.notFull.Wait()
	}

	queue.items[queue.head] = e
	queue.head++
	if queue.head == len(queue.items) {
		queue.head = 0
	}
	queue.count++
	queue.notEmpty.Signal()
}

func (queue *ArrayBlockingQueue[T]) Dequeue() (e T) {

	for queue.count == 0 {
		queue.notEmpty.Wait()
	}

	e = queue.items[queue.tail]
	queue.items[queue.tail] = nil
	queue.tail++
	if queue.tail == len(queue.items) {
		queue.tail = 0
	}
	queue.count--
	queue.notFull.Signal()
	return e
}
