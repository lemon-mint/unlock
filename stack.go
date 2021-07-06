package unlock

import (
	"sync/atomic"
	"unsafe"
)

// Stack: Treiber's stack
type Stack struct {
	top *node
	len int64
}

func NewStack() *Stack {
	return new(Stack)
}

func (s *Stack) Push(x unsafe.Pointer) {
	newHead := new(node)
	newHead.value = x
	for {
		oldHead := s.top
		newHead.next = oldHead
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&s.top)),
			unsafe.Pointer(oldHead),
			unsafe.Pointer(newHead),
		) {
			break
		}
	}
	atomic.AddInt64(&s.len, 1)
}

func (s *Stack) Pop() (unsafe.Pointer, bool) {
	var x unsafe.Pointer
	for {
		oldHead := s.top
		if oldHead == nil {
			return nil, false
		}
		newHead := oldHead.next
		x = oldHead.value
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&s.top)),
			unsafe.Pointer(oldHead),
			unsafe.Pointer(newHead),
		) {
			break
		}
	}
	atomic.AddInt64(&s.len, -1)
	return x, true
}

func (s Stack) Len() int64 {
	return s.len
}
