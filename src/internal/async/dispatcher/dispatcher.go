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
	channel      chan status.CompletionStatus
	sem          *semaphore.Semaphore
	conn         abstract.IDBConnection
	requestsRepo repository.AsyncRequestRepository
	//workerFactory *worker.WorkerFactory
	sleepTime   time.Duration
	maxAttempts int
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
		sem:          sem,
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

				requests, err := d.requestsRepo.GetAllRequests(ctx, txConn)
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
				if err = d.requestsRepo.SetStatusAndIncrementById(ctx, txConn, requst.ID, status.StatusInProgress); err != nil {
					txConn.Rollback()
					d.sem.Release()
					return
				}

				if err := txConn.Commit(); err != nil {
					d.sem.Release()
					return
				}
				// ToDO: сделать фабрику worker и запустить worker

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
			if completionStatus.Err != nil {
				if completionStatus.Attempts >= d.maxAttempts {
					if err := d.requestsRepo.SetStatusById(ctx, d.conn, completionStatus.RequestId, status.StatusFailed); err != nil {
						log.Printf("can not set status for request %d: %v", completionStatus.RequestId, err)

					} else {
						if err := d.requestsRepo.SetStatusAndIncrementById(ctx, d.conn, completionStatus.RequestId, status.StatusQueued); err != nil {
							log.Printf("can not set status for request %d: %v", completionStatus.RequestId, err)
						}
					}

				} else {
					if err := d.requestsRepo.SetStatusById(ctx, d.conn, completionStatus.RequestId, status.StatusSuccess); err != nil {
						log.Printf("can not set status for request %d: %v", completionStatus.RequestId, err)
					}
				}
			}
		default:
			if err := d.requestsRepo.CloseTimeoutRequests(ctx, d.conn); err != nil {
				log.Printf("error closing timeout requests: %v", err)
			}
			time.Sleep(d.sleepTime)
		}
	}
}
