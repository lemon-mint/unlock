package unlock

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type RingBuffer struct {
	buf      []unsafe.Pointer
	padding0 [8]uint64
	size     int
	padding1 [8]uint64
	r        int
	w        int
	padding2 [8]uint64
	counter  int64
	padding3 [8]uint64
	TLock
}

func (b *RingBuffer) dummy() {
	_ = b.padding0
	_ = b.padding1
	_ = b.padding2
	_ = b.padding3
}

func NewRingBuffer(size int) *RingBuffer {
	r := new(RingBuffer)
	r.buf = make([]unsafe.Pointer, size)
	r.size = size
	r.dummy()
	return r
}

func (b *RingBuffer) EnQueue(x unsafe.Pointer) {
	for {
		ctr := atomic.LoadInt64(&b.counter)
		if ctr+1 > int64(b.size) {
			runtime.Gosched()
			continue
		}
		if atomic.CompareAndSwapInt64(&b.counter, ctr, ctr+1) {
			break
		}
	}
	b.Lock()
	b.buf[b.w] = x
	b.w++
	if b.w >= b.size {
		b.w = 0
	}
	b.Unlock()
}

func (b *RingBuffer) DeQueue() unsafe.Pointer {
	for {
		ctr := atomic.LoadInt64(&b.counter)
		if ctr <= 0 {
			runtime.Gosched()
			continue
		}
		if atomic.CompareAndSwapInt64(&b.counter, ctr, ctr-1) {
			break
		}
	}
	b.Lock()
	val := b.buf[b.r]
	b.r++
	if b.r >= b.size {
		b.r = 0
	}
	b.Unlock()
	return val
}

func (b *RingBuffer) EnQueueMany(x []unsafe.Pointer) {
	length := len(x)
	for {
		ctr := atomic.LoadInt64(&b.counter)
		if ctr+int64(length) > int64(b.size) {
			runtime.Gosched()
			continue
		}
		if atomic.CompareAndSwapInt64(&b.counter, ctr, ctr+int64(length)) {
			break
		}
	}
	b.Lock()
	for i := range x {
		b.buf[b.w] = x[i]
		b.w++
		if b.w >= b.size {
			b.w = 0
		}
	}
	b.Unlock()
}

func (b *RingBuffer) DeQueueMany(dst []unsafe.Pointer) {
	length := len(dst)
	for {
		ctr := atomic.LoadInt64(&b.counter)
		if ctr < int64(length) {
			runtime.Gosched()
			continue
		}
		if atomic.CompareAndSwapInt64(&b.counter, ctr, ctr-int64(length)) {
			break
		}
	}
	b.Lock()
	for i := range dst {
		dst[i] = b.buf[b.r]
		b.r++
		if b.r >= b.size {
			b.r = 0
		}
	}
	b.Unlock()
}
