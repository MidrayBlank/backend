package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type MaleFemaleAgeExtractor struct {
	RosstatCSVReader[dao.MaleFemaleAgeExtracted]
	oktmoColumn  int
	yearColumn   int
	sexColumn    int
	ageColumn    int
	amountColumn int
}

func NewMaleFemaleAgeExtractor() *MaleFemaleAgeExtractor {
	return &MaleFemaleAgeExtractor{
		oktmoColumn:  11,
		yearColumn:   18,
		sexColumn:    4,
		ageColumn:    5,
		amountColumn: 19,
	}
}

func (ext *MaleFemaleAgeExtractor) extractRow(row []string) *dao.MaleFemaleAgeExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	sexStr := strings.TrimSpace(row[ext.sexColumn])
	ageStr := strings.TrimSpace(row[ext.ageColumn])
	amountStr := strings.TrimSpace(row[ext.amountColumn])

	maleAmount := 0
	femaleAmount := 0

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	age, err := strconv.Atoi(ageStr)

	if err != nil {
		return nil
	}

	amount, err := strconv.Atoi(strings.ReplaceAll(amountStr, ".0", ""))

	if err != nil {
		return nil
	}

	switch sexStr {
	case "Мужчины":
		maleAmount = amount
	case "Женщины":
		femaleAmount = amount
	default:
		return nil
	}

	return &dao.MaleFemaleAgeExtracted{
		Code:         oktmo,
		Year:         year,
		Age:          age,
		MaleAmount:   maleAmount,
		FemaleAmount: femaleAmount,
	}
}

func (ext *MaleFemaleAgeExtractor) Extract(ctx context.Context, filePath string) ([]*dao.MaleFemaleAgeExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
