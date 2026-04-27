package semaphore

type Semaphore struct {
	channel chan struct{}
}

<<<<<<< HEAD
func NewSemaphore(size int) Semaphore {
	return Semaphore{
=======
func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
>>>>>>> a89c928 (POS-53: Init Dispatcher)
		channel: make(chan struct{}, size),
	}
}

<<<<<<< HEAD
func (s Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s Semaphore) Release() {
=======
func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s *Semaphore) Release() {
>>>>>>> a89c928 (POS-53: Init Dispatcher)
	<-s.channel
}
