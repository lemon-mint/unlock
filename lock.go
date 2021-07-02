package unlock

import (
	"runtime"
	"sync/atomic"
)

type Lock struct {
	lock uint32
}

// Lock acquires the lock
func (l *Lock) Lock() {
	for {
		if l.TryLock() {
			break // escape loop
		} else {
			runtime.Gosched() // yeilds to other goroutines
		}
	}
}

// TryLock tries to acquire the lock
func (l *Lock) TryLock() (locked bool) {
	return atomic.CompareAndSwapUint32(&l.lock, 0, 1)
}

// Unlock releases the lock
func (l *Lock) Unlock() {
	atomic.StoreUint32(&l.lock, 0) // release lock
}
