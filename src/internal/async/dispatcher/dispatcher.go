package dispatcher

import (
	"context"
	"log"
	"time"

	request_map "backend/src/internal/async/request_map/impl"
	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/db/abstract"

	worker_config "backend/src/internal/async/worker/config"
	repository "backend/src/internal/repository/impl"
)

type Dispatcher struct {
	channel      chan status.CompletionStatus
	sem          semaphore.Semaphore
	conn         abstract.IDBConnection
	requestsRepo repository.AsyncRequestRepository

	//workerFactory *worker.WorkerFactory
	workerConfig worker_config.WorkerConfig

	sleepTime          time.Duration
	maxAttempts        int
	requestAttemptsMap request_map.RequestAttemptsMap
}

func NewDispatcher(
	conn abstract.IDBConnection,
	requestsRepo repository.AsyncRequestRepository,
	workerConfig worker_config.WorkerConfig,
	sleepTime time.Duration,
	maxAttempts int,
	maxWorkersCount int,
) *Dispatcher {
	sem := semaphore.NewSemaphore(maxWorkersCount)
	return &Dispatcher{
		channel:            make(chan status.CompletionStatus, maxWorkersCount*2),
		sem:                sem,
		conn:               conn,
		requestsRepo:       requestsRepo,
		workerConfig:       workerConfig,
		sleepTime:          sleepTime,
		requestAttemptsMap: request_map.NewRequestAttemptsMap(),
		maxAttempts:        maxAttempts,
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

			d.fetchRequestsAndRunWorkers(ctx)
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

				attempts, exists := d.requestAttemptsMap.Get(completionStatus.RequestId)
				if !exists {
					log.Printf("request %d not found", completionStatus.RequestId)
					continue
				}

				if attempts >= d.maxAttempts {
					if err := d.requestsRepo.SetStatusById(ctx, d.conn, completionStatus.RequestId, status.StatusFailed); err != nil {
						log.Printf("can not set status for request %d: %v", completionStatus.RequestId, err)

					}
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

		default:
			if err := d.requestsRepo.CloseTimeoutRequests(ctx, d.conn); err != nil {
				log.Printf("error closing timeout requests: %v", err)
			}
			time.Sleep(d.sleepTime)
		}
	}
}

func (d *Dispatcher) fetchRequestsAndRunWorkers(ctx context.Context) {
	txConn := d.conn.BeginTx()

	defer func() {
		if r := recover(); r != nil {
			txConn.Rollback()
			d.sem.Release()
		}
	}()

	request, err := d.requestsRepo.GetOneRequest(ctx, txConn)
	if err != nil {
		txConn.Rollback()
		d.sem.Release()
		return
	}

	if request == nil {
		txConn.Rollback()
		d.sem.Release()
		time.Sleep(d.sleepTime)
		return
	}

	d.requestAttemptsMap.Put(request.ID, request.Attempts)

	if err = d.requestsRepo.SetStatusAndIncrementById(ctx, txConn, request.ID, status.StatusInProgress); err != nil {
		txConn.Rollback()
		d.sem.Release()
		return
	}

	if err := txConn.Commit(); err != nil {
		d.sem.Release()
		return
	}

	// ToDO: сделать фабрику worker и запустить worker
}
