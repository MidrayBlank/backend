package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RosstatParser struct{}

func NewRosstatParser() *RosstatParser {
	return &RosstatParser{}
}

func (p *RosstatParser) ParseRosstat(filePath string) ([]RosstatRawRecord, error) {
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

	isHousing := strings.Contains(filePath, "Y48010001")
	isSalary := strings.Contains(filePath, "Y48423007")

	var records []RosstatRawRecord

	for _, row := range rows[1:] {
		if len(row) < 18 {
			continue
		}

		var oktmo string
		var year int
		var value float64

		if isHousing || isSalary {
			if len(row) > 10 {
				oktmo = strings.TrimSpace(row[10])
			}
			if len(row) > 17 {
				year, _ = strconv.Atoi(strings.TrimSpace(row[17]))
			}
			if len(row) > 18 {
				value = parseFloat(row[18])
			}
		} else {
			if len(row) > 9 {
				oktmo = strings.TrimSpace(row[9])
			}
			if len(row) > 16 {
				year, _ = strconv.Atoi(strings.TrimSpace(row[16]))
			}
			if len(row) > 17 {
				value = parseFloat(row[17])
			}
		}

		if oktmo == "" || year == 0 {
			continue
		}

		records = append(records, RosstatRawRecord{
			Oktmo: oktmo,
			Year:  year,
			Value: value,
		})
	}

	return records, nil
}
