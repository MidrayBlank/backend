package parser

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"backend/pkg/parser/rosstat/downloader"
)

type LandParser struct{}

func NewLandParser() *LandParser {
	return &LandParser{}
}

func (p *LandParser) Parse(ctx context.Context, filePath string) ([]downloader.RosstatRawRecord, error) {
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

	var records []downloader.RosstatRawRecord

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

		oktmo := strings.TrimSpace(row[9])
		year, _ := strconv.Atoi(strings.TrimSpace(row[16]))
		value := parseFloat(row[17])

		if oktmo == "" || year == 0 {
			continue
		}

		records = append(records, downloader.RosstatRawRecord{
			Oktmo: oktmo,
			Year:  year,
			Value: value,
		})
	}

	return records, nil
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
