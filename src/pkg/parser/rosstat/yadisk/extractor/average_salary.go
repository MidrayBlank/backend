package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type AverageSalaryExtractor struct {
	RosstatCSVReader[dao.AverageSalaryExtracted]
	oktmoColumn         int
	yearColumn          int
	averageSalaryColumn int
}

func NewAverageSalaryExtractor() *AverageSalaryExtractor {
	return &AverageSalaryExtractor{
		oktmoColumn:         10,
		yearColumn:          17,
		averageSalaryColumn: 18,
	}
}

func (ext *AverageSalaryExtractor) extractRow(row []string) *dao.AverageSalaryExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	averageSalaryStr := strings.TrimSpace(row[ext.averageSalaryColumn])

	oktmo, err := strconv.Atoi(oktmoStr)
	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	averageSalary, err := strconv.ParseFloat(averageSalaryStr, 64)

	if err != nil {
		return nil
	}

	return &dao.AverageSalaryExtracted{
		Code:          oktmo,
		Year:          year,
		AverageSalary: averageSalary,
	}
}

func (ext *AverageSalaryExtractor) Extract(ctx context.Context, filePath string) ([]*dao.AverageSalaryExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
