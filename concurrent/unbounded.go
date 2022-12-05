package concurrent

import (
	"sync"
)

/**** YOU CANNOT MODIFY ANY OF THE FOLLOWING INTERFACES/TYPES ********/
type Task interface{}

type DEQueue interface {
	PushBottom(task Task)
	IsEmpty() bool //returns whether the queue is empty
	PopTop() Task
	PopBottom() Task
	Size() int // Allowed modification
}

/******** DO NOT MODIFY ANY OF THE ABOVE INTERFACES/TYPES *********************/

type node struct {
	task Task
	prev *node
	next *node
}

type dequeue struct {
	mtx  sync.Mutex
	head *node
	tail *node
	size int
}

// NewUnBoundedDEQueue returns an empty UnBoundedDEQueue
func NewUnBoundedDEQueue() DEQueue {
	dq := &dequeue{
		mtx:  sync.Mutex{},
		head: &node{},
		tail: &node{},
		size: 0,
	}
	dq.head.next = dq.tail
	dq.tail.prev = dq.head
	return dq
}

func (dq *dequeue) PushBottom(task Task) {
	dq.mtx.Lock()
	// enq(dq, task)
	newNode := node{task: task}
	newNode.next = dq.tail
	newNode.prev = dq.tail.prev
	newNode.prev.next = &newNode
	dq.tail.prev = &newNode
	dq.size += 1
	dq.mtx.Unlock()
}

func (dq *dequeue) IsEmpty() bool {
	if dq.head.next == dq.tail {
		return true
	} else {
		return false
	}
}

func (dq *dequeue) PopTop() Task {
	dq.mtx.Lock()
	defer dq.mtx.Unlock()
	if dq.IsEmpty() {
		return nil
	}
	task := dq.head.next.task
	dq.head.next = dq.head.next.next
	dq.head.next.prev = dq.head
	dq.size -= 1
	return task
}

func (dq *dequeue) PopBottom() Task {
	dq.mtx.Lock()
	defer dq.mtx.Unlock()
	if dq.IsEmpty() {
		return nil
	}
	task := dq.tail.prev.task
	dq.tail.prev = dq.tail.prev.prev
	dq.tail.prev.next = dq.tail
	dq.size -= 1
	return task
}

func (dq *dequeue) Size() int {
	dq.mtx.Lock()
	defer dq.mtx.Unlock()
	return dq.size
}

// Balance two queues
func balancing(q0 DEQueue, q1 DEQueue, threshold int) {
	var qMin DEQueue
	var qMax DEQueue

	if q0.Size() < q1.Size() {
		qMin = q0
		qMax = q1
	} else {
		qMin = q1
		qMax = q0
	}
	// Get diff
	diff := (qMax.Size() - qMin.Size()) / 2
	if diff > threshold {
		// Move the task until balanced
		for qMax.Size() > qMin.Size() {
			qMin.PushBottom(qMax.PopTop())
		}
	}
}
