package worker

import (
	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/db/abstract"
)

type WorkerFactory struct {
	handlers map[int]WorkerFuncType
	conn     abstract.IDBConnection
	repos    repository.WorkerRepositories
	ch       chan status.CompletionStatus
	sem      *semaphore.Semaphore
}
