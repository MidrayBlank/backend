package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type SchoolsExtractor struct {
	RosstatCSVReader[dao.SchoolsExtracted]
	oktmoColumn   int
	yearColumn    int
	schoolsColumn int
}

func NewSchoolsExtractor() *SchoolsExtractor {
	return &SchoolsExtractor{
		oktmoColumn:   9,
		yearColumn:    16,
		schoolsColumn: 17,
	}
}

func (ext *SchoolsExtractor) extractRow(row []string) *dao.SchoolsExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	schoolsStr := strings.TrimSpace(row[ext.schoolsColumn])

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	schools, err := strconv.Atoi(strings.ReplaceAll(schoolsStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.SchoolsExtracted{
		Code:    oktmo,
		Year:    year,
		Schools: schools,
	}
}

func (ext *SchoolsExtractor) Extract(ctx context.Context, filePath string) ([]*dao.SchoolsExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
