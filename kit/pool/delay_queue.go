package pool

import "container/heap"

type PriorityQueue struct {
	items []*Item
}

type Item struct {
	WorkerTask
	priority int64
	index    int
}

func (queue *PriorityQueue) NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items: make([]*Item, 0, 10),
	}
}

func (queue *PriorityQueue) Len() int {
	return len(queue.items)
}

func (queue *PriorityQueue) Less(i, j int) bool {
	return queue.items[i].priority > queue.items[j].priority
}

func (queue *PriorityQueue) Swap(i, j int) {
	queue.items[i], queue.items[j] = queue.items[j], queue.items[i]
	queue.items[i].index = i
	queue.items[j].index = j
}

func (queue *PriorityQueue) Push(x any) {
	n := queue.Len()
	item := x.(*Item)
	item.index = n
	(*queue).items = append((*queue).items, item)
}

func (queue *PriorityQueue) Pop() any {
	old := *queue
	n := len(old.items)
	item := old.items[n-1]
	old.items[n-1] = nil // avoid memory leak
	item.index = -1      // for safety
	(*queue).items = old.items[0 : n-1]
	return item
}

func (queue *PriorityQueue) update(item *Item, value WorkerTask, priority int64) {
	item.WorkerTask = value
	item.priority = priority
	heap.Fix(queue, item.index)
}
