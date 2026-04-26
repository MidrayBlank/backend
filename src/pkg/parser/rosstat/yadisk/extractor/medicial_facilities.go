package extractor

import (
	"context"
	"strconv"
	"strings"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type MedicialFacilitiesExtractor struct {
	RosstatCSVReader[dao.MedicialFacilitiesExtracted]
	oktmoColumn             int
	yearColumn              int
	medicalFacilitiesColumn int
}

func NewMedicialFacilitiesExtractor() *MedicialFacilitiesExtractor {
	return &MedicialFacilitiesExtractor{
		oktmoColumn:             9,
		yearColumn:              16,
		medicalFacilitiesColumn: 17,
	}
}

func (ext *MedicialFacilitiesExtractor) extractRow(row []string) *dao.MedicialFacilitiesExtracted {
	oktmoStr := strings.TrimSpace(row[ext.oktmoColumn])
	yearStr := strings.TrimSpace(row[ext.yearColumn])
	medicalFacilitiesStr := strings.TrimSpace(row[ext.medicalFacilitiesColumn])

	oktmo, err := strconv.Atoi(oktmoStr)

	if err != nil {
		return nil
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return nil
	}

	medicalFacilities, err := strconv.Atoi(strings.ReplaceAll(medicalFacilitiesStr, ".0", ""))

	if err != nil {
		return nil
	}

	return &dao.MedicialFacilitiesExtracted{
		Code:               oktmo,
		Year:               year,
		MedicialFacilities: medicalFacilities,
	}
}

func (ext *MedicialFacilitiesExtractor) Extract(ctx context.Context, filePath string) ([]*dao.MedicialFacilitiesExtracted, error) {
	return ext.ExtractRows(ctx, filePath, ext.extractRow)
}
