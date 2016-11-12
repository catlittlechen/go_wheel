package channel

import (
	"sync"
)

// Element 队列的点
type Element struct {
	next  *Element
	value string
}

// Queue 队列
type Queue struct {
	count  int
	push   *sync.Mutex
	pop    *sync.Mutex
	null   *sync.Mutex
	change *sync.Mutex
	head   *Element
	tail   *Element
}

// Init 初始化队列
func (queue *Queue) Init() {
	queue.push = new(sync.Mutex)
	queue.pop = new(sync.Mutex)
	queue.null = new(sync.Mutex)
	queue.change = new(sync.Mutex)
	queue.null.Lock()
}

// Push 推入队列
func (queue *Queue) Push(value string) {

	queue.push.Lock()
	defer queue.push.Unlock()

	e := &Element{
		value: value,
	}

	queue.change.Lock()
	queue.count++
	if queue.tail == nil {
		queue.head = e
		queue.tail = e
		queue.null.Unlock()
	} else {
		queue.tail.next = e
		queue.tail = e
	}
	queue.change.Unlock()

	return
}

// Pop 取出队列
func (queue *Queue) Pop() (value string) {

	queue.pop.Lock()
	defer queue.pop.Unlock()

	queue.null.Lock()
	value = queue.head.value

	queue.change.Lock()
	queue.count--
	queue.head = queue.head.next
	if queue.head == nil {
		queue.tail = nil
	} else {
		queue.null.Unlock()
	}
	queue.change.Unlock()

	return
}
