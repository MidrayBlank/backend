package impl

import "sync"

type RequestAttemptsMap struct {
	data map[int]int
	mu   sync.Mutex
}

func NewRequestAttemptsMap() RequestAttemptsMap {
	return RequestAttemptsMap{
		data: make(map[int]int),
	}
}

func (ram *RequestAttemptsMap) Put(requestID int) int {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	ram.data[requestID]++
	return ram.data[requestID]
}

func (ram *RequestAttemptsMap) Get(requestID int) int {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	return ram.data[requestID]
}

func (ram *RequestAttemptsMap) Has(requestID int) bool {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	_, exists := ram.data[requestID]
	return exists
}
