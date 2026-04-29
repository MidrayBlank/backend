package ask_ai

import (
	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/async/worker/repository"
	"backend/src/internal/config"
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	ai "backend/src/pkg/ai/openrouter"
	"context"
	"fmt"
	"log"
	"strings"
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

	stats, err := prepareStatistics(ctx, conn, repositories, regionCode)
	if err != nil {
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	regionName := getRegionName(conn, repositories, regionCode)

	tableStats := formatStatsForAI(stats)
	prompt := buildAIPrompt(regionCode, regionName, tableStats)

	cfg := config.Load()
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
) string {
	geo, err := repositories.GeoRepository.GetGeoByCode(conn, code)
	if err != nil {
		return fmt.Sprintf("region with code %d is not founded", code)
	}
	return geo.Name
}

func prepareStatistics(
	ctx context.Context,
	conn abstract.IDBConnection,
	repositories repository.WorkerRepositories,
	code int,
) ([]*domain.Rosstat, error) {
	rosstat, err := repositories.RosstatRepository.GetRosstatByCodeForLastFiveYears(conn, code)
	if err != nil {
		return nil, err
	}
	return rosstat, nil
}

func buildAIPrompt(code int, regionName string, stats string) string {
	prompt := fmt.Sprintf(`Ты профессиональный демограф-аналитик.

Составь аналитическую справку по демографической ситуации в регионе "%s" (код %d).

## Все данные по региону(за последние 5 лет):
%s

## ТРЕБОВАНИЯ К ОТВЕТУ:
Ответ должен быть ТОЛЬКО в формате JSON, без лишнего текста до или после.

Структура ответа:
{
  "summary": "краткое резюме (3-4 предложения с ключевыми цифрами)",
  "forecast": {
    "scenario": "базовый сценарий (2025-2030)",
    "population": "ожидаемая численность населения (диапазон)",
    "natural_growth": "естественный прирост в процентах",
    "migration_growth": "миграционный прирост"
  },
  "trends": [
    "тенденция 1",
    "тенденция 2",
    "тенденция 3"
  ],
  "recommendations": [
    "рекомендация 1",
    "рекомендация 2",
    "рекомендация 3"
  ],
}

## ВАЖНО:
- summary должен содержать конкретные цифры и проценты
- trends должно быть 3-5 пунктов
- recommendations должно быть 3-5 пунктов
- confidence — число от 0 до 1
- Используй ТОЛЬКО русский язык
- Никакого текста вне JSON!

Ответь ТОЛЬКО JSON.`, regionName, code, stats)

	return prompt
}

func formatStatsForAI(stats []*domain.Rosstat) string {
	var sb strings.Builder
	sb.WriteString("| Год | Население | Рождаемость | Смертность | Прибывшие | Уехавшие  |\n")
	sb.WriteString("|-----|-----------|-------------|------------|-----------|-----------|\n")

	for _, s := range stats {
		population := formatIntValue(s.PopulationAmount)
		birth := formatIntValue(s.BirthAmount)
		death := formatIntValue(s.DeathAmount)

		arrival := formatIntValue(s.ArrivalAmount)
		departure := formatIntValue(s.DepartureAmount)

		sb.WriteString(fmt.Sprintf("\t| %d \t| %s \t| %s \t| %s \t| %s \t| %s \t|\n",
			s.Year, population, birth, death, arrival, departure))
	}

	return sb.String()
}

func formatIntValue(value *int) string {
	if value == nil {
		return "—"
	}
	return fmt.Sprintf("%d", *value)
}
