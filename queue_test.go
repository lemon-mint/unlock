package unlock

import (
	"testing"
	"unsafe"
)

func BenchmarkEnQueueDeQueue(b *testing.B) {
	q := NewQueue()
	b.RunParallel(func(p *testing.PB) {
		var ctr int
		for p.Next() {
			q.EnQueue(unsafe.Pointer(&ctr))
			q.DeQueue()
		}
	})
}

func BenchmarkEnQueue(b *testing.B) {
	q := NewQueue()
	b.RunParallel(func(p *testing.PB) {
		var ctr int
		for p.Next() {
			q.EnQueue(unsafe.Pointer(&ctr))
		}
	})
}
