package unlock

import (
	"testing"
	"unsafe"
)

func BenchmarkRingBufferEnQueueDeQueue(b *testing.B) {
	r := NewRingBuffer(2048)
	b.RunParallel(func(p *testing.PB) {
		var ctr int
		for p.Next() {
			r.EnQueue(unsafe.Pointer(&ctr))
			r.DeQueue()
		}
	})
}

func BenchmarkRingBufferMany(b *testing.B) {
	r := NewRingBuffer(2048)
	b.RunParallel(func(p *testing.PB) {
		var buf [8]unsafe.Pointer

		for p.Next() {
			r.EnQueueMany(buf[:])
			r.DeQueueMany(buf[:])
		}
	})
}
