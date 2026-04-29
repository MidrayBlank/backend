package yadisk_parse_test

import (
	"context"
	"testing"

	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/async/worker/config"
	"backend/src/internal/async/worker/yadisk_parse"
	"backend/src/internal/db/postgres"
	"backend/src/internal/repository/impl"
)

func TestPopulation(t *testing.T) {
	ctx := context.Background()
	ch := make(chan status.CompletionStatus)
	sem := semaphore.NewSemaphore(3)
	conn := postgres.NewPostgresConnection("postgresql://test_user:123@127.0.0.1:5432/test_db?sslmode=disable")
	repositories := config.WorkerConfig{
		GeoRepository:        impl.NewGeoRepository(),
		RosstatRepository:    impl.NewRosstatRepository(),
		RosstatAgeRepository: impl.NewRosstatAgeRepository(),
		AiApiRepository:      impl.NewAiApiRepository(),
	}

	yadisk_parse.YadiskParseWorker(
		ctx,
		ch,
		sem,
		conn,
		repositories,
		1,
		0,
	)

	t.Log("YadiskWorker was finished.\n")
}
