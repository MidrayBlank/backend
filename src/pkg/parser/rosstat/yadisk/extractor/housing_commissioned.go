package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type HousingCommissionedExtractor struct {
	RosstatCSVReader[dao.HousingCommissionedExtracted]
	oktmoColumn               int
	yearColumn                int
	housingCommissionedColumn int
}

func NewHousingCommissionedExtractor() *HousingCommissionedExtractor {
	return &HousingCommissionedExtractor{
		oktmoColumn:               10,
		yearColumn:                17,
		housingCommissionedColumn: 18,
	}
}

func (ext *HousingCommissionedExtractor) extractRow(row []string) *dao.HousingCommissionedExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	housingCommissionedStr := strings.TrimSpace(row[ext.housingCommissionedColumn])

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	housingCommissioned, err := strconv.Atoi(strings.ReplaceAll(housingCommissionedStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.HousingCommissionedExtracted{
		Code:                oktmo,
		Year:                year,
		HousingCommissioned: housingCommissioned,
	}
}

func (ext *HousingCommissionedExtractor) Extract(ctx context.Context, filePath string) ([]*dao.HousingCommissionedExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
