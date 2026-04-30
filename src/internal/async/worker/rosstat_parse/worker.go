package rosstat_parse

import (
	"context"

	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/async/worker/repository"
	"backend/src/internal/db/abstract"
)

const batchSize = 2000

func RosstatParseWorker(
	ctx context.Context,
	ch chan status.CompletionStatus,
	sem semaphore.Semaphore,
	conn abstract.IDBConnection,
	repositories repository.WorkerRepositories,
	requestId int,
	param int,
) {

}
