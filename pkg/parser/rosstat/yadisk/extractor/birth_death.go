package extractor

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BirthDeathParser struct{}

func NewBirthDeathParser() *BirthDeathParser {
	return &BirthDeathParser{}
}

func (p *BirthDeathParser) Parse(ctx context.Context, filePath string) ([]OfficialRecord, error) {
	// ПРОВЕРКА 1: перед открытием файла
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

	// ПРОВЕРКА 2: после открытия файла
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// ПРОВЕРКА 3: после чтения файла
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows")
	}

	var records []OfficialRecord

	for i, row := range rows[1:] {
		// ПРОВЕРКА 4: каждые 10000 строк
		if i%10000 == 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}

		if len(row) < 18 {
			continue
		}

		oktmo := strings.TrimSpace(row[9])
		if oktmo == "" || oktmo == "CD" || len(oktmo) < 8 {
			continue
		}

		year, _ := strconv.Atoi(strings.TrimSpace(row[16]))
		if year == 0 {
			continue
		}

		valueStr := strings.TrimSpace(row[17])
		valueStr = strings.ReplaceAll(valueStr, " ", "")
		valueStr = strings.ReplaceAll(valueStr, "\u00A0", "")

		value, _ := strconv.ParseFloat(valueStr, 64)
		if value == 0 {
			continue
		}

		records = append(records, OfficialRecord{
			Oktmo: oktmo,
			Year:  year,
			Value: value,
		})
	}

	return records, nil
}
