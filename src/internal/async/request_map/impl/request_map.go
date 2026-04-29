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

func (ram *RequestAttemptsMap) Put(requestID int, attempts int) {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	ram.data[requestID] = attempts + 1
}

func (ram *RequestAttemptsMap) Get(requestID int) (int, bool) {
	ram.mu.Lock()
	defer ram.mu.Unlock()
	value, exists := ram.data[requestID]

	return value, exists
}

func (ram *RequestAttemptsMap) Delete(requestID int) {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	delete(ram.data, requestID)
}
