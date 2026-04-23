package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DemographyParser struct{}

func NewDemographyParser() *DemographyParser {
	return &DemographyParser{}
}

func (p *DemographyParser) ParseDemography(filePath string) ([]DemographyRawRecord, error) {
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

	var records []DemographyRawRecord

	for _, row := range rows[1:] {
		if len(row) < 14 {
			continue
		}

		record := DemographyRawRecord{
			Oktmo:             strings.TrimSpace(row[0]),
			NotZato:           parseInt(row[1]),
			Region:            strings.TrimSpace(row[2]),
			MunType:           strings.TrimSpace(row[3]),
			Municipality:      strings.TrimSpace(row[4]),
			Year:              parseInt(row[5]),
			Population:        parseInt(row[6]),
			AveragePopulation: parseFloat(row[7]),
			Deaths:            parseInt(row[8]),
			Births:            parseInt(row[9]),
			Migration:         parseInt(row[10]),
			MortalityRate:     parseFloat(row[11]),
			BirthRate:         parseFloat(row[12]),
			MigrationRate:     parseFloat(row[13]),
		}

		records = append(records, record)
	}

	return records, nil
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
