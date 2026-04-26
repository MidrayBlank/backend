package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type ArrivalExtractor struct {
	RosstatCSVReader[dao.ArrivalExtracted]
	oktmoColumn         int
	yearColumn          int
	arrivalColumn       int
	migrationTypeColumn int
	sexTypeColumn       int
	ageTypeColumn       int
}

func NewArrivalExtractor() *ArrivalExtractor {
	return &ArrivalExtractor{
		oktmoColumn:         12,
		yearColumn:          19,
		arrivalColumn:       20,
		migrationTypeColumn: 4,
		sexTypeColumn:       5,
		ageTypeColumn:       6,
	}
}

func (ext *ArrivalExtractor) extractRow(row []string) *dao.ArrivalExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	arrivalStr := strings.TrimSpace(row[ext.arrivalColumn])
	migrationType := strings.TrimSpace(row[ext.migrationTypeColumn])
	sexType := strings.TrimSpace(row[ext.sexTypeColumn])
	ageType := strings.TrimSpace(row[ext.ageTypeColumn])

	if !strings.Contains(migrationType, "всего") || !strings.Contains(sexType, "Всего") || !strings.Contains(ageType, "Всего") {
		return nil
	}

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	arrival, err := strconv.Atoi(strings.ReplaceAll(arrivalStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.ArrivalExtracted{
		Code:    oktmo,
		Year:    year,
		Arrival: arrival,
	}
}

func (ext *ArrivalExtractor) Extract(ctx context.Context, filePath string) ([]*dao.ArrivalExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
