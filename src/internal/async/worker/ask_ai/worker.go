package ask_ai

import (
	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	worker_config "backend/src/internal/async/worker/config"
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
	workerConfig worker_config.WorkerConfig,
	requestId int,
	regionCode int,
) {
	defer sem.Release()

	rosstatStats, rosstatAgeStats, err := prepareStatistics(ctx, conn, workerConfig, regionCode, workerConfig.AiReportYears)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	regionName, err := getRegionName(conn, workerConfig, regionCode)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	tableStats := formatStatsForAI(rosstatStats, rosstatAgeStats)
	prompt := buildAIPrompt(regionName, tableStats)

	aiRouterClient := ai.NewAIOpenRouter(workerConfig.AiApiKeys, workerConfig.AiModel, workerConfig.OpenRouterBaseURL)

	response, err := aiRouterClient.SendRequest(prompt)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	err = workerConfig.AiReportRepository.Upsert(conn, regionCode, response)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	err = workerConfig.AiApiRepository.IncreaseRequests(conn, workerConfig.ApiKeys)
	if err != nil {
		log.Printf("Failed to increase count of requests: %v", err)
	}

	ch <- status.CompletionStatus{RequestId: requestId, Err: nil}
}

func getRegionName(
	conn abstract.IDBConnection,
	workerConfig worker_config.WorkerConfig,
	code int,
) (string, error) {
	geo, err := workerConfig.GeoRepository.GetGeoByCode(conn, code)
	if err != nil {
		return "", err
	}
	return geo.Name, nil
}
