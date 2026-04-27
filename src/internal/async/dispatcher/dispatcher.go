package dispatcher

import (
	"context"
	"log"
	"time"

	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/db/abstract"
	repository "backend/src/internal/repository/impl"
)

type Dispatcher struct {
	channel       chan status.CompletionStatus
	sem           *semaphore.Semaphore
	conn          abstract.IDBConnection
	requestsRepo  repository.AsyncRequestRepository
	workerFactory *worker.WorkerFactory
	sleepTime     time.Duration
	maxAttempts   int
}

func NewDispatcher(
	conn abstract.IDBConnection,
	requestsRepo repository.AsyncRequestRepository,
	sleepTime time.Duration,
	maxAttempts int,
	maxWorkersCount int,
) *Dispatcher {
	sem := semaphore.NewSemaphore(maxWorkersCount)
	return &Dispatcher{
		channel:      make(chan status.CompletionStatus, maxWorkersCount*2),
		sem:          &sem,
		conn:         conn,
		requestsRepo: requestsRepo,
		sleepTime:    sleepTime,
		maxAttempts:  maxAttempts,
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

			func() {
				txConn := d.conn.BeginTx()

				defer func() {
					if r := recover(); r != nil {
						txConn.Rollback()
						d.sem.Release()
					}
				}()

				requests, err := d.requestsRepo.GetAllRequests(ctx, d.conn)
				if err != nil {
					txConn.Rollback()
					d.sem.Release()
					return
				}

				if len(requests) == 0 {
					txConn.Rollback()
					d.sem.Release()
					time.Sleep(d.sleepTime)
					return
				}

				requst := requests[0]
				if err = d.requestsRepo.SetStatusAndIncrementById(ctx, d.conn, requst.ID, status.StatusInProgress); err != nil {
					txConn.Rollback()
					d.sem.Release()
					return
				}

				if err := txConn.Commit(); err != nil {
					d.sem.Release()
					return
				}
				// ToDO: сделать фабрику worker и запусть worker

			}()
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
			if err := d.requestsRepo.CloseTimeoutRequests(ctx, d.conn); err != nil {
				log.Printf("Error closing timeout requests: %v", err)
			}
			time.Sleep(d.sleepTime)
		}
	}
}
