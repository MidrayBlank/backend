package report

import "backend/src/pkg/report/population_growth/domain"

type PopulationGrowthCalculator struct{}

func NewPopulationGrowthCalculator() *PopulationGrowthCalculator {
	return &PopulationGrowthCalculator{}
}

func (c *PopulationGrowthCalculator) CompileReport(data []*domain.PopulationGrowthParams) []*domain.PopulationGrowthReport {
	if len(data) == 0 {
		return []*domain.PopulationGrowthReport{}
	}
	report := make([]*domain.PopulationGrowthReport, 0, len(data))
	for _, dat := range data {
		if dat.PopulationAmountFrom == 0 {
			continue
		}
		growthPct := float32(dat.PopulationAmountTo-dat.PopulationAmountFrom) / float32(dat.PopulationAmountFrom) * 100
		report = append(report, &domain.PopulationGrowthReport{
			OKTMO:               dat.OKTMO,
			YearFrom:            dat.YearFrom,
			YearTo:              dat.YearTo,
			PopulationGrowthPct: growthPct,
		})
	}
	return report
}
