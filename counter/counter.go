package counter

import (
	"sync"
)

type TotalCounter struct {
	Total int
	mutex sync.Mutex
}

func NewCounter() (tc *TotalCounter) {
	tc = new(TotalCounter)
	return

}

func (tc *TotalCounter) SafeAdd(count int) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	tc.Total += count

}
