package dispatcher

import (
	"context"
	"log"

	"backend/src/internal/async/semaphore"
	"backend/src/internal/db/abstract"
)

type CompletionStatus struct {
	RequestId int
	Err       error
}

// добавить sleepTime
type Dispatcher struct {
	channel chan CompletionStatus
	sem     *semaphore.Semaphore

	conn abstract.IDBConnection
}

func NewDispatcher(
	maxWorkersCount int,
) *Dispatcher {
	return &Dispatcher{
		channel: make(chan CompletionStatus, maxWorkersCount*2),
		sem:     semaphore.NewSemaphore(maxWorkersCount),
	}
}

func (d *Dispatcher) RunWorkers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("RunWorkers finished: %s", ctx.Err().Error())
			return
		default:
			d.sem.Acquire()

			// <запрашиваем одну заявку под SELECT FOR UPDATE в транзакции>
			// НАШЛИ ЗАЯВКУ
			// создаем горутину с параметрами из заявки - по типу заявки создаем воркер (фабрика)
			// завершаем транзакцию
			// НЕ НАШЛИ ЗАЯВКУ
			// уходим в сон на некоторое время
		}
	}
}

func (d *Dispatcher) ProcessCompletion(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("ProcessCompletion finished: %s", ctx.Err().Error())
			return
		case completionStatus := <-d.channel:
			// Помечаем заявку как
			// - SUCCESS, если успешно завершилась
			// - FAILED, если завершилась с ошибкой, а attempts == maxAttempts
			// - QUEUED, если завершилась с ошибкой, а attempts < maxAttempts
			// <помечаем заявки IN_PROGRESS, у которых вышел таймаут и больше 3-ех ошибок как FAILED> (?)
		default:
			// уходим в сон на некоторое время
		}
	}
}
