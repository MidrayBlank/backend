package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type DeathExtractor struct {
	RosstatCSVReader[dao.DeathExtracted]
	oktmoColumn int
	yearColumn  int
	deathColumn int
}

func NewDeathExtractor() *DeathExtractor {
	return &DeathExtractor{
		oktmoColumn: 9,
		yearColumn:  16,
		deathColumn: 17,
	}
}

func (ext *DeathExtractor) extractRow(row []string) *dao.DeathExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	deathStr := strings.TrimSpace(row[ext.deathColumn])

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	death, err := strconv.Atoi(strings.ReplaceAll(deathStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.DeathExtracted{
		Code:  oktmo,
		Year:  year,
		Death: death,
	}
}

func (ext *DeathExtractor) Extract(ctx context.Context, filePath string) ([]*dao.DeathExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
