package parser

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AgeSexParser struct{}

func NewAgeSexParser() *AgeSexParser {
	return &AgeSexParser{}
}

func (p *AgeSexParser) ParseAgeSex(ctx context.Context, filePath string) ([]AgeSexRawRecord, error) {
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

	var records []AgeSexRawRecord

	for i, row := range rows[1:] {
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

		sex := strings.TrimSpace(row[4])
		if sex != "Мужчины" && sex != "Женщины" {
			continue
		}
		age := strings.TrimSpace(row[6])

		oktmo := strings.TrimSpace(row[11])
		year, _ := strconv.Atoi(strings.TrimSpace(row[17]))
		value, _ := strconv.Atoi(strings.TrimSpace(row[18]))

		if oktmo == "" || year == 0 || age == "" {
			continue
		}

		records = append(records, AgeSexRawRecord{
			Oktmo: oktmo,
			Year:  year,
			Age:   age,
			Sex:   sex,
			Value: value,
		})
	}

	return records, nil
}

// AgeSexRawRecord - сырые данные половозрастной структуры
type AgeSexRawRecord struct {
	Oktmo string
	Year  int
	Age   string
	Sex   string
	Value int
}
