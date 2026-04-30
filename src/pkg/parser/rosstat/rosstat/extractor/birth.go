package extractor

import (
	"context"
	"strconv"
	"strings"

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
	println("=== NewBirthExtractor created ===")
	return &BirthExtractor{storage: storage}
}

func (ext *BirthExtractor) Extract(ctx context.Context, filePath string) error {
	println("=== BirthExtractor.Extract START ===")
	println("File path:", filePath)

	manager := state_manager.NewStateManager[dao.BirthExtracted](
		state_manager.StateMapFuncType[dao.BirthExtracted]{
			SKIP:              ext.skip,
			GET_YEARS:         ext.getYears,
			GET_NAME_OR_BIRTH: ext.getNameOrBirth,
		},
	)

	err := ext.ExtractRows(ctx, filePath, manager)
	println("=== BirthExtractor.Extract END, error:", err)
	return err
}

func (ext *BirthExtractor) skip(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	println("DEBUG skip: row[0] =", row[0])
	manager.ChangeState(GET_YEARS)
}

func (ext *BirthExtractor) getYears(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	println("DEBUG getYears: processing row, length:", len(row))

	ext.years = make([]int, len(row)-1)

	for i := 1; i < len(row); i++ {
		if row[i] == "" {
			continue
		}
		year, err := strconv.Atoi(row[i])
		if err == nil {
			ext.years[i-1] = year
			println("DEBUG getYears: year", i, "=", year)
		}
	}

	println("DEBUG getYears: total years parsed:", len(ext.years))
	manager.ChangeState(GET_NAME_OR_BIRTH)
}

func (ext *BirthExtractor) getNameOrBirth(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	if len(row) == 0 || row[0] == "" {
		println("DEBUG getNameOrBirth: empty row, skipping")
		return
	}

	println("DEBUG getNameOrBirth: row[0] =", row[0])
	println("DEBUG getNameOrBirth: current ext.code =", ext.code)

	// Проверяем, есть ли числа в строке
	hasNumbers := false
	for i := 1; i < len(row) && i < 5; i++ {
		if row[i] != "" && row[i] != ";" {
			hasNumbers = true
			break
		}
	}
	println("DEBUG getNameOrBirth: hasNumbers =", hasNumbers)

	if strings.HasPrefix(row[0], "Муниципальный район") && hasNumbers {
		println("DEBUG getNameOrBirth: MATCH! Calling getBirth with code", ext.code)
		ext.getBirth(row, manager)
		return
	}

	code := ext.storage.GetCode(row[0])
	if code != 0 {
		println("DEBUG getNameOrBirth: Found code", code, "for name:", row[0])
		ext.code = code
	} else {
		println("DEBUG getNameOrBirth: No code found for name:", row[0])
	}
}

func (ext *BirthExtractor) getBirth(row []string, manager *state_manager.StateManager[dao.BirthExtracted]) {
	println("DEBUG getBirth: called with code", ext.code, "years len", len(ext.years))

	birthDAOs := make([]*dao.BirthExtracted, 0, len(ext.years))

	for i := 1; i < len(row) && i-1 < len(ext.years); i++ {
		if row[i] == "" {
			println("DEBUG getBirth: row[", i, "] is empty")
			continue
		}
		birth, err := strconv.Atoi(row[i])
		if err == nil && birth > 0 {
			println("DEBUG getBirth: year", ext.years[i-1], "birth", birth)
			birthDAOs = append(birthDAOs, &dao.BirthExtracted{
				Code:  ext.code,
				Year:  ext.years[i-1],
				Birth: birth,
			})
		} else {
			println("DEBUG getBirth: failed to parse row[", i, "] =", row[i], "error:", err)
		}
	}

	if len(birthDAOs) > 0 {
		ext.storage.SetBirth(birthDAOs)
		println("Saved birth data for code", ext.code, "years count", len(birthDAOs))
	} else {
		println("No birth data for code", ext.code)
	}
}
