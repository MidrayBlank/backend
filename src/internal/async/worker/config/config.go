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
}

func NewWorkerConfig() WorkerConfig {
	return WorkerConfig{
		GeoRepository:        impl.NewGeoRepository(),
		RosstatRepository:    impl.NewRosstatRepository(),
		RosstatAgeRepository: impl.NewRosstatAgeRepository(),
		AiApiRepository:      impl.NewAiApiRepository(),
		AiReportRepository:   impl.NewAiReportRepository(),
	}
}
