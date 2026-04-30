package main

import (
	"context"

	"backend/src/internal/config"
	"backend/src/internal/db/postgres"
	"backend/src/internal/domain"
	"backend/src/internal/repository/impl"
	"backend/src/pkg/parser/rosstat/rosstat"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	config := config.Load()
	conn := postgres.NewPostgresConnection(config.GetDBDSN())

	geoRepo := impl.NewGeoRepository()
	rosstatRepo := impl.NewRosstatRepository()

	parser := rosstat.NewRosstatParser()

	ctx := context.Background()
	parsedSlice, err := parser.Parse(ctx)

	if err != nil {
		panic(err)
	}

	higherGEOs := make([]*domain.Geo, len(parsedSlice))
	lowerGEOs := make([]*domain.Geo, len(parsedSlice))
	rosstats := make([]*domain.Rosstat, len(parsedSlice))

	for index, parsed := range parsedSlice {
		if parsed.ParentCode < 100 {
			higherGEOs = append(higherGEOs, &domain.Geo{
				Code:       parsed.Code,
				ParentCode: &parsed.ParentCode,
				Name:       parsed.Name,
				Level:      2,
			})
		} else {
			lowerGEOs = append(lowerGEOs, &domain.Geo{
				Code:       parsed.Code,
				ParentCode: &parsed.ParentCode,
				Name:       parsed.Name,
				Level:      3,
			})
		}

		rosstats[index] = &domain.Rosstat{
			Code: parsed.Code,
			Year: parsed.Year,
		}
	}

	geoRepo.UpsertBatch(conn, higherGEOs)
	geoRepo.UpsertBatch(conn, lowerGEOs)

	rosstatRepo.UpsertBatch(conn, rosstats)
}
