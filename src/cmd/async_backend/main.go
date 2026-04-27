package async_backend

import (
	"backend/src/internal/async/dispatcher"
	"backend/src/internal/config"
	"backend/src/internal/db/postgres"
	rimpl "backend/src/internal/repository/impl"
	"context"
	"sync"
	"time"
)

func main() {
	config := config.Load()

	conn := postgres.NewPostgresConnection(config.GetDBDSN())
	requestsRepo := rimpl.NewAsyncRequesRepository()
	dispatcher := dispatcher.NewDispatcher(conn, *requestsRepo, 5*time.Second, config.MaxAttempts, config.MaxWorkersCount)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx := context.Background()

	go dispatcher.RunWorkers(ctx)
	go dispatcher.ProcessCompletion(ctx)

	wg.Wait()
}
