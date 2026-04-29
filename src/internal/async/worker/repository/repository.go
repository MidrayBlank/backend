package repository

import (
	"backend/src/internal/repository/abstract"
	"backend/src/internal/repository/impl"
)

type WorkerRepositories struct {
	GeoRepository        abstract.IGeoRepository
	RosstatRepository    abstract.IRosstatRepository
	RosstatAgeRepository abstract.IRosstatAgeRepository

	AiApiRepository abstract.IAiApiRepository
}

func NewWorkerRepositories() WorkerRepositories {
	return WorkerRepositories{
		GeoRepository:        impl.NewGeoRepository(),
		RosstatRepository:    impl.NewRosstatRepository(),
		RosstatAgeRepository: impl.NewRosstatAgeRepository(),
		AiApiRepository:      impl.NewAiApiRepository(),
	}
}
