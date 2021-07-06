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
