package semaphore

type Semaphore struct {
	channel chan struct{}
}

func NewSemaphore(size int) Semaphore {
	return Semaphore{
		channel: make(chan struct{}, size),
	}
}

func (s Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s Semaphore) Release() {
	<-s.channel
}
