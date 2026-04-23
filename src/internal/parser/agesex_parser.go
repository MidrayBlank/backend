package parser

import (
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

func (p *AgeSexParser) ParseAgeSex(filePath string) ([]AgeSexRawRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows")
	}

	var records []AgeSexRawRecord

	for _, row := range rows[1:] {
		if len(row) < 18 {
			continue
		}

		sex := strings.TrimSpace(row[4])
		if sex != "Мужчины" && sex != "Женщины" {
			continue
		}

		oktmo := strings.TrimSpace(row[11])
		year, _ := strconv.Atoi(strings.TrimSpace(row[17]))
		value, _ := strconv.Atoi(strings.TrimSpace(row[18]))

		if oktmo == "" || year == 0 {
			continue
		}

		records = append(records, AgeSexRawRecord{
			Oktmo: oktmo,
			Year:  year,
			Sex:   sex,
			Value: value,
		})
	}

	return records, nil
}
