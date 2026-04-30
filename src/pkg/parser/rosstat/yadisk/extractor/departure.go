package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type DepartureExtractor struct {
	RosstatCSVReader[dao.DepartureExtracted]
	oktmoColumn         int
	yearColumn          int
	departureColumn     int
	migrationTypeColumn int
	sexTypeColumn       int
	ageTypeColumn       int
}

func NewDepartureExtractor() *DepartureExtractor {
	return &DepartureExtractor{
		oktmoColumn:         12,
		yearColumn:          19,
		departureColumn:     20,
		migrationTypeColumn: 4,
		sexTypeColumn:       5,
		ageTypeColumn:       6,
	}
}

func (ext *DepartureExtractor) extractRow(row []string) *dao.DepartureExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	departureStr := strings.TrimSpace(row[ext.departureColumn])
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

	departure, err := strconv.Atoi(strings.ReplaceAll(departureStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.DepartureExtracted{
		Code:      oktmo,
		Year:      year,
		Departure: departure,
	}
}

func (ext *DepartureExtractor) Extract(ctx context.Context, filePath string) ([]*dao.DepartureExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
