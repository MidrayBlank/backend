package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type LandAreaExtractor struct {
	RosstatCSVReader[dao.LandAreaExtracted]
	oktmoColumn   int
	yearColumn    int
	landAreColumn int
}

func NewLandAreaExtractor() *LandAreaExtractor {
	return &LandAreaExtractor{
		oktmoColumn:   9,
		yearColumn:    16,
		landAreColumn: 17,
	}
}

func (ext *LandAreaExtractor) extractRow(row []string) *dao.LandAreaExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	landAreaStr := strings.TrimSpace(row[ext.landAreColumn])

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	landArea, err := strconv.Atoi(strings.ReplaceAll(landAreaStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.LandAreaExtracted{
		Code:     oktmo,
		Year:     year,
		LandArea: landArea,
	}
}

func (ext *LandAreaExtractor) Extract(ctx context.Context, filePath string) ([]*dao.LandAreaExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
