package repository

import "backend/src/internal/repository/abstract"

type WorkerRepositories struct {
	GeoRepository        abstract.IGeoRepository
	RosstatRepository    abstract.IRosstatRepository
	RosstatAgeRepository abstract.IRosstatAgeRepository

	AiApiRepository abstract.IAiApiRepository
}
