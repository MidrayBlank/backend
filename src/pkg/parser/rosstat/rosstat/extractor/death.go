package extractor

import (
	"context"
	"strconv"

	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"backend/src/pkg/parser/rosstat/rosstat/state_manager"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
)

type DeathExtractor struct {
	BaseCSVExtractor[dao.DeathExtracted]

	storage *storage.Storage

	years []int
	code  int
}

func NewDeathExtractor(storage *storage.Storage) *DeathExtractor {
	return &DeathExtractor{storage: storage}
}

func (ext *DeathExtractor) Extract(ctx context.Context, filePath string) error {
	manager := state_manager.NewStateManager[dao.DeathExtracted](
		state_manager.StateMapFuncType[dao.DeathExtracted]{
			SKIP:              ext.skip,
			GET_YEARS:         ext.getYears,
			GET_NAME_OR_DEATH: ext.getNameOrDeath,
		},
	)

	return ext.ExtractRows(ctx, filePath, manager)
}

func (ext *DeathExtractor) skip(row []string, manager *state_manager.StateManager[dao.DeathExtracted]) {
	manager.ChangeState(GET_YEARS)
}

func (ext *DeathExtractor) getYears(row []string, manager *state_manager.StateManager[dao.DeathExtracted]) {
	ext.years = make([]int, len(row)-1)

	for i := 1; i < len(row); i++ {
		if row[i] == "" {
			continue
		}
		year, err := strconv.Atoi(row[i])
		if err == nil {
			ext.years[i-1] = year
		}
	}

	manager.ChangeState(GET_NAME_OR_DEATH)
}

func (ext *DeathExtractor) getNameOrDeath(row []string, manager *state_manager.StateManager[dao.DeathExtracted]) {
	if len(row) == 0 || row[0] == "" {
		return
	}

	// Проверяем, является ли строка строкой с данными
	if (row[0] == "Муниципальный район") || row[0] == "Муниципальный округ" ||
		row[0] == "Городской округ, городской округ с внутригородским делением" ||
		row[0] == "Городские поселения" || row[0] == "Сельские поселения" && len(row) > 1 && row[1] != "" {
		ext.getDeath(row, manager)
		return
	}

	code := ext.storage.GetCode(row[0])
	if code != 0 {
		ext.code = code
	}
}

func (ext *DeathExtractor) getDeath(row []string, manager *state_manager.StateManager[dao.DeathExtracted]) {
	deathDAOs := make([]*dao.DeathExtracted, 0, len(ext.years))

	for i := 1; i < len(row) && i-1 < len(ext.years); i++ {
		if row[i] == "" {
			continue
		}
		death, err := strconv.Atoi(row[i])
		if err == nil && death > 0 {
			deathDAOs = append(deathDAOs, &dao.DeathExtracted{
				Code:  ext.code,
				Year:  ext.years[i-1],
				Death: death,
			})
		}
	}

	if len(deathDAOs) > 0 {
		ext.storage.SetDeath(deathDAOs)
	}
}
