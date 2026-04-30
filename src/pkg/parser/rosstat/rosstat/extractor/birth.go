package extractor

import (
	"context"
	"strconv"

	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"backend/src/pkg/parser/rosstat/rosstat/state_manager"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
)

type BirthExtractor struct {
	BaseCSVExtractor[dao.BirthExtracted]

	storage *storage.Storage

	years []int
	code  int
}

func NewBirthExtractor(storage *storage.Storage) *BirthExtractor {
	return &BirthExtractor{storage: storage}
}

func (ext *BirthExtractor) Extract(ctx context.Context, filePath string) error {
	manager := state_manager.NewStateManager[dao.BirthExtracted](
		state_manager.StateMapFuncType[dao.BirthExtracted]{
			SKIP:              ext.skip,
			GET_YEARS:         ext.getYears,
			GET_NAME_OR_BIRTH: ext.getNameOrBirth,
		},
	)

	return ext.ExtractRows(ctx, filePath, manager)
}

func (ext *BirthExtractor) skip(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	manager.ChangeState(GET_YEARS)
}

func (ext *BirthExtractor) getYears(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	ext.years = make([]int, len(row)-1)

	for i := 1; i < len(ext.years); i++ {
		ext.years[i], _ = strconv.Atoi(row[i])
	}

	manager.ChangeState(GET_NAME_OR_POPULATION)
}

func (ext *BirthExtractor) getNameOrBirth(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	if len(row) > 0 && row[0] == "Число родившихся (без мертворожденных), человек, значение показателя за год" {
		ext.getBirth(row, manager)
		return
	}

	code := ext.storage.GetCode(row[0])
	if code != 0 {
		ext.code = code
	}
}

func (ext *BirthExtractor) getBirth(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	birthDAOs := make([]*dao.BirthExtracted, 0, len(ext.years))

	for i := 1; i < len(row) && i-1 < len(ext.years); i++ {
		if row[i] == "" {
			continue
		}

		birth, err := strconv.Atoi(row[i])
		if err == nil {
			birthDAOs = append(birthDAOs, &dao.BirthExtracted{
				Code:  ext.code,
				Year:  ext.years[i-1],
				Birth: birth,
			})
		}
	}

	if len(birthDAOs) > 0 {
		ext.storage.SetBirth(birthDAOs)
	}
}
