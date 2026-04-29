package config

import (
	"backend/src/internal/repository/abstract"
	"backend/src/internal/repository/impl"
)

type WorkerConfig struct {
	GeoRepository        abstract.IGeoRepository
	RosstatRepository    abstract.IRosstatRepository
	RosstatAgeRepository abstract.IRosstatAgeRepository

	AiApiRepository    abstract.IAiApiRepository
	AiReportRepository abstract.IAiReportRepository

	ApiKeys           []string
	AiModel           string
	OpenRouterBaseURL string
	AiReportYears     int
}

func NewWorkerConfig(
	apiKeys []string,
	aiModel string,
	openRoterBaseURL string,
	aiReportYears int,
) WorkerConfig {
	return WorkerConfig{
		GeoRepository:        impl.NewGeoRepository(),
		RosstatRepository:    impl.NewRosstatRepository(),
		RosstatAgeRepository: impl.NewRosstatAgeRepository(),
		AiApiRepository:      impl.NewAiApiRepository(),
		AiReportRepository:   impl.NewAiReportRepository(),

		ApiKeys:           apiKeys,
		AiModel:           aiModel,
		OpenRouterBaseURL: openRoterBaseURL,
		AiReportYears:     aiReportYears,
	}
}
