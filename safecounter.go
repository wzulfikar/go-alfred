package alfred

import "sync"

type SafeCounter struct {
	sync.Mutex
	Count int
}

func (counter *SafeCounter) Increment(n int) {
	counter.Lock()
	counter.Count += n
	counter.Unlock()
}

func (counter *SafeCounter) Decrement(n int) {
	counter.Lock()
	counter.Count -= n
	counter.Unlock()
}
