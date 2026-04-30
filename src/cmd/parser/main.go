package main

import (
	"context"
	"fmt"
	"log"

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

	fmt.Printf("Parsed count: %d\n", len(parsedSlice))

	GEOs := make([]*domain.Geo, 0, len(parsedSlice))
	rosstats := make([]*domain.Rosstat, len(parsedSlice))

	codeMap := make(map[int]bool)

	fmt.Printf("Grouping...\n")

	for index, parsed := range parsedSlice {
		if _, exists := codeMap[parsed.Code]; !exists {
			codeMap[parsed.Code] = true
			GEOs = append(GEOs, &domain.Geo{
				Code:       parsed.Code,
				ParentCode: &parsed.ParentCode,
				Name:       parsed.Name,
				Level:      2,
			})
		}

		rosstats[index] = &domain.Rosstat{
			Code:             parsed.Code,
			Year:             parsed.Year,
			PopulationAmount: parsed.Population,
		}
	}

	fmt.Printf("Grouped\n")

	batchSize := 50
	startIndex := 0
	endIndex := batchSize

	fmt.Printf("Upsert batch geo higher...\n")
	for {
		if endIndex >= len(GEOs) {
			endIndex = len(GEOs)
		}
		fmt.Printf("Upsert batch geo higher [%d:%d]...\n", startIndex, endIndex)
		if err := geoRepo.UpsertBatch(conn, GEOs[startIndex:endIndex]); err != nil {
			log.Printf("ERROR [higherGEOs]: %s\n", err.Error())
			return
		}

		if endIndex >= len(GEOs) {
			break
		}

		startIndex = endIndex
		endIndex += batchSize

		if endIndex >= len(GEOs) {
			endIndex = len(GEOs)
		}
	}

	startIndex = 0
	endIndex = batchSize

	fmt.Printf("Upsert batch rosstat...\n")
	for {
		if endIndex >= len(rosstats) {
			endIndex = len(rosstats)
		}
		fmt.Printf("Upsert batch rosstat [%d:%d]...\n", startIndex, endIndex)
		if err := rosstatRepo.UpsertBatch(conn, rosstats[startIndex:endIndex]); err != nil {
			log.Printf("ERROR [lowerGEOs]: %s\n", err.Error())
			return
		}

		if endIndex >= len(rosstats) {
			break
		}

		startIndex = endIndex
		endIndex += batchSize

		if endIndex >= len(rosstats) {
			endIndex = len(rosstats)
		}
	}
}

func getSubjectCode(code int) int {
	extendedSubjectCode := code / 100000
	subjectCode := code / 1000000

	switch extendedSubjectCode {
	case 118:
		subjectCode = extendedSubjectCode
	case 718:
		subjectCode = extendedSubjectCode
	case 719:
		subjectCode = extendedSubjectCode
	}

	return subjectCode
}
