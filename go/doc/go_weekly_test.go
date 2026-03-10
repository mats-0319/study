package doc

import (
	"sync"
	"testing"
)

type Locked[T any] struct {
	mu sync.Mutex
	v  T
}

func NewLocked[T any](value T) *Locked[T] {
	return &Locked[T]{v: value}
}

func (l *Locked[T]) Get() T {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.v
}

func (l *Locked[T]) Set(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.v = value
}

func TestGoWeekly(t *testing.T) {
	for range 5 {
		counter := NewLocked(0)

		wg := sync.WaitGroup{}

		for range 10 {
			wg.Go(func() {
				for range 1000 {
					v := counter.Get()
					v++
					counter.Set(v)
				}
			})
		}

		wg.Wait()

		res := counter.Get()
		t.Log(res)
	}
}
