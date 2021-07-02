package unlock

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	value unsafe.Pointer
	next  *node
}

type Queue struct {
	head *node
	tail *node
}

func (q *Queue) EnQueue(x unsafe.Pointer) bool {
	newNode := new(node)
	newNode.value = x
	var tail *node
	for {
		//tail = (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail))))
		//next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
		tail = q.tail
		next := tail.next
		if next != nil {
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(next))
			continue
		}
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)), nil, unsafe.Pointer(newNode)) {
			break
		}
	}
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(newNode))
}

func (q *Queue) DeQueue() (unsafe.Pointer, bool) {
	for {
		firstNode := q.head
		lastNode := q.tail
		nextNode := firstNode.next
		if firstNode != q.head {
			continue
		}
		if firstNode == lastNode {
			if nextNode == nil {
				return nil, false
			}
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(lastNode), unsafe.Pointer(nextNode))
			continue
		}
		x := nextNode.value
		if !atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(firstNode), unsafe.Pointer(nextNode)) {
			continue
		}
		return x, true
	}
}

func NewQueue() *Queue {
	n := new(node)
	q := new(Queue)
	q.head = n
	q.tail = n
	return q
}
