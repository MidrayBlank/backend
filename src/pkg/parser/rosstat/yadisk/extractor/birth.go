package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type BirthExtractor struct {
	RosstatCSVReader[dao.BirthExtracted]
	oktmoColumn int
	yearColumn  int
	birthColumn int
}

func NewBirthExtractor() *BirthExtractor {
	return &BirthExtractor{
		oktmoColumn: 9,
		yearColumn:  16,
		birthColumn: 17,
	}
}

func (ext *BirthExtractor) extractRow(row []string) *dao.BirthExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	birthStr := strings.TrimSpace(row[ext.birthColumn])

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	birth, err := strconv.Atoi(strings.ReplaceAll(birthStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.BirthExtracted{
		Code:  oktmo,
		Year:  year,
		Birth: birth,
	}
}

func (ext *BirthExtractor) Extract(ctx context.Context, filePath string) ([]*dao.BirthExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
