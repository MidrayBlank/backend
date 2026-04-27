package async_backend

import (
	"backend/src/internal/async/dispatcher"
	"backend/src/internal/config"
	"context"
	"sync"
)

func main() {
	config := config.Load()
	dispatcher := dispatcher.NewDispatcher(config.MaxWorkersCount)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx := context.Background()

	go dispatcher.RunWorkers(ctx)
	go dispatcher.ProcessCompletion(ctx)

	wg.Wait()
}
