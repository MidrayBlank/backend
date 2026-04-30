package extractor

import (
	"context"
	"strconv"

	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"backend/src/pkg/parser/rosstat/rosstat/state_manager"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
)

type PopulationExtractor struct {
	BaseCSVExtractor[dao.PopulationExtracted]

	storage *storage.Storage

	years []int
	code  int
}

func NewPopulationExtractor(storage *storage.Storage) *PopulationExtractor {
	return &PopulationExtractor{storage: storage}
}

func (ext *PopulationExtractor) Extract(ctx context.Context, filePath string) error {
	manager := state_manager.NewStateManager[dao.PopulationExtracted](
		state_manager.StateMapFuncType[dao.PopulationExtracted]{
			SKIP:                   ext.skip,
			GET_YEARS:              ext.getYears,
			GET_NAME_OR_POPULATION: ext.getNameOrPopulation,
		},
	)

	return ext.ExtractRows(ctx, filePath, manager)
}

const (
	SKIP = iota
	GET_YEARS
	GET_NAME_OR_POPULATION
	GET_NAME_OR_BIRTH
	GET_NAME_OR_DEATH
)

func (ext *PopulationExtractor) skip(row []string, manager *state_manager.StateManager[dao.PopulationExtracted]) {
	manager.ChangeState(GET_YEARS)
}

func (ext *PopulationExtractor) getYears(row []string, manager *state_manager.StateManager[dao.PopulationExtracted]) {
	ext.years = make([]int, len(row)-1)

	for i := 1; i < len(ext.years); i++ {
		ext.years[i], _ = strconv.Atoi(row[i])
	}

	manager.ChangeState(GET_NAME_OR_POPULATION)
}

func (ext *PopulationExtractor) getNameOrPopulation(row []string, manager *state_manager.StateManager[dao.PopulationExtracted]) {
	if row[0] == "Все население" {
		ext.getPopulation(row, manager)
	}

	code := ext.storage.GetCode(row[0])

	if code != 0 {
		ext.code = ext.storage.GetCode(row[0])
	}
}

func (ext *PopulationExtractor) getPopulation(row []string, manager *state_manager.StateManager[dao.PopulationExtracted]) {
	populationDAOs := make([]*dao.PopulationExtracted, 0, len(ext.years)-1)

	maxIndex := min(len(ext.years), len(row))

	for i := 1; i < maxIndex; i++ {
		population, err := strconv.Atoi(row[i])

		if err == nil {
			populationDAOs = append(populationDAOs, &dao.PopulationExtracted{
				Code:       ext.code,
				Year:       ext.years[i],
				Population: population,
			})
		}
	}

	ext.storage.SetPopulation(populationDAOs)
}
