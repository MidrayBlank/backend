package ask_ai

import (
	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/async/worker/repository"
	"backend/src/internal/config"
	"backend/src/internal/db/abstract"
	ai "backend/src/pkg/ai/openrouter"
	"context"
	"log"
)

func AskAIByCodeWorker(
	ctx context.Context,
	ch chan status.CompletionStatus,
	sem semaphore.Semaphore,
	conn abstract.IDBConnection,
	repositories repository.WorkerRepositories,
	requestId int,
	regionCode int,
) {
	defer sem.Release()
	cfg := config.Load()

	rosstatStats, rosstatAgeStats, err := prepareStatistics(ctx, conn, repositories, regionCode, cfg.AIReportYears)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	regionName, err := getRegionName(conn, repositories, regionCode)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	tableStats := formatStatsForAI(rosstatStats, rosstatAgeStats)
	prompt := buildAIPrompt(regionName, tableStats)

	aiRouterClient := ai.NewAIOpenRouter(cfg.AiApiKey, cfg.AiModel, cfg.AiBaseUrl)

	response, err := aiRouterClient.SendRequest(prompt)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	err = repositories.AiReportRepository.Upsert(conn, regionCode, response)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	err = repositories.AiApiRepository.IncreaseRequests(conn, cfg.AiApiKey)
	if err != nil {
		log.Printf("Failed to increase count of requests: %v", err)
	}

	ch <- status.CompletionStatus{RequestId: requestId, Err: nil}
}

func getRegionName(
	conn abstract.IDBConnection,
	repositories repository.WorkerRepositories,
	code int,
) (string, error) {
	geo, err := repositories.GeoRepository.GetGeoByCode(conn, code)
	if err != nil {
		return "", err
	}
	return geo.Name, nil
}
