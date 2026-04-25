package extractor

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"backend/pkg/parser/rosstat/yadisk/downloader"
)

type MigrationParser struct{}

func NewMigrationParser() *MigrationParser {
	return &MigrationParser{}
}

func (p *MigrationParser) ParseMigration(ctx context.Context, filePath string) ([]downloader.MigrationRecord, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows")
	}

	const (
		OKTMO_COL  = 12
		YEAR_COL   = 19
		VALUE_COL  = 20
		GRUP_2_COL = 5
		VOZR_COL   = 6
	)

	var records []downloader.MigrationRecord

	for i, row := range rows[1:] {
		if i%10000 == 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}

		if len(row) <= VALUE_COL {
			continue
		}

		grup2 := strings.TrimSpace(row[GRUP_2_COL])
		vozr := strings.TrimSpace(row[VOZR_COL])

		if grup2 != "Всего" || vozr != "Всего" {
			continue
		}

		oktmo := strings.TrimSpace(row[OKTMO_COL])
		if oktmo == "" {
			continue
		}

		if strings.HasSuffix(oktmo, "000") {
			continue
		}

		year, err := strconv.Atoi(strings.TrimSpace(row[YEAR_COL]))
		if err != nil {
			continue
		}

		valueStr := strings.TrimSpace(row[VALUE_COL])
		valueStr = strings.ReplaceAll(valueStr, " ", "")
		valueStr = strings.ReplaceAll(valueStr, "\u00A0", "")

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		records = append(records, downloader.MigrationRecord{
			Oktmo: oktmo,
			Year:  year,
			Value: int(value),
		})
	}

	fmt.Printf("  %s: %d records\n", filepath.Base(filePath), len(records))
	return records, nil
}
