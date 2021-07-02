package unlock

import (
	"runtime"
	"sync"
	"testing"
)

func TestLock_Lock(t *testing.T) {
	var counter int
	var l = new(Lock)
	var wg sync.WaitGroup
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func() {
			runtime.Gosched()
			for j := 0; j < 2048; j++ {
				l.Lock()
				counter++
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	t.Run("lock counter", func(t *testing.T) {
		if counter != 256*2048 {
			t.Error("race condition is occurred")
		}
	})
}

var lock Lock
var mu sync.Mutex

func BenchmarkSpinLock(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var counter int
		for p.Next() {
			lock.Lock()
			counter++
			lock.Unlock()
		}
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var counter int
		for p.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}
