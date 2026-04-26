package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type PopulationExtractor struct {
	RosstatCSVReader[dao.PopulationExtracted]
	oktmoColumn          int
	yearColumn           int
	populationColumn     int
	populationTypeColumn int
}

func NewPopulationExtractor() *PopulationExtractor {
	return &PopulationExtractor{
		oktmoColumn:          10,
		yearColumn:           17,
		populationColumn:     18,
		populationTypeColumn: 4,
	}
}

func (ext *PopulationExtractor) extractRow(row []string) *dao.PopulationExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	populationStr := strings.TrimSpace(row[ext.populationColumn])
	populationType := strings.TrimSpace(row[ext.populationTypeColumn])

	urbanUsed := false
	ruralUsed := false

	switch populationType {
	case "Все население":
		urbanUsed = true
		ruralUsed = true
	case "Городское население":
		urbanUsed = true
	case "Сельское население":
		ruralUsed = true
	}

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	population, err := strconv.Atoi(strings.ReplaceAll(populationStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.PopulationExtracted{
		Code:       oktmo,
		Year:       year,
		Population: population,
		UrbanUsed:  urbanUsed,
		RuralUsed:  ruralUsed,
	}
}

func (ext *PopulationExtractor) Extract(ctx context.Context, filePath string) ([]*dao.PopulationExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
