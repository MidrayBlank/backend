package abstract

import (
	"backend/src/internal/domain"
)

type IReportService interface {
	GetGrowthReport(codes []int, yearFrom, yearTo int) ([]domain.GrowthReport, error)
	GetDensityReport(codes []int, year int) ([]domain.DensityReport, error)
	GetGrowthLeadersReport(codes []int, yearFrom, yearTo, limit int) (*domain.GrowthLeadersResponse, error)
	GetIndicatorReport(codes []int, year int) ([]domain.IndicatorsReport, error)
}
