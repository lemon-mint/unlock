package unlock

import (
	"testing"
	"unsafe"
)

func BenchmarkPushStack(b *testing.B) {
	q := NewStack()
	b.RunParallel(func(p *testing.PB) {
		var ctr int
		for p.Next() {
			q.Push(unsafe.Pointer(&ctr))
		}
	})
}

func TestStack_Pop(t *testing.T) {
	stack := NewStack()
	stack.Push(nil)
	stack.Push(nil)
	stack.Push(nil)
	stack.Push(nil)
	stack.Push(nil)
	var ctr = 0
	t.Run("Get Value", func(t *testing.T) {
		ctr++
		if _, ok := stack.Pop(); !ok {
			if ctr != 5 {
				t.Errorf("Pop() = %v expect %v", ctr, 5)
			}
		}
	})

}
