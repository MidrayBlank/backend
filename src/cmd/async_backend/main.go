package async_backend

import (
	"backend/src/internal/async/dispatcher"
	worker_repositories "backend/src/internal/async/worker/repository"
	"backend/src/internal/config"
	"backend/src/internal/db/postgres"
	rimpl "backend/src/internal/repository/impl"
	"context"
	"sync"
)

func main() {
	config := config.Load()

	conn := postgres.NewPostgresConnection(config.GetDBDSN())
	requestsRepo := rimpl.NewAsyncRequesRepository()
	workerRepositories := worker_repositories.NewWorkerRepositories()

	dispatcher := dispatcher.NewDispatcher(conn, *requestsRepo, workerRepositories, config.DispatcherTimeSleep, config.DispatcherMaxAttempts, config.MaxWorkersCount)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx := context.Background()

	go dispatcher.RunWorkers(ctx)
	go dispatcher.ProcessCompletion(ctx)

	wg.Wait()
}
