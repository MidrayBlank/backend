package report

import "backend/src/pkg/report/sync/population_growth/domain"

type PopulationGrowthReportCompiler struct{}

func NewPopulationGrowthReportCompiler() *PopulationGrowthReportCompiler {
	return &PopulationGrowthReportCompiler{}
}

func (c *PopulationGrowthReportCompiler) Compile(data []*domain.PopulationGrowthParams) []*domain.PopulationGrowthReport {
	if len(data) == 0 {
		return []*domain.PopulationGrowthReport{}
	}
	report := make([]*domain.PopulationGrowthReport, 0, len(data))
	for _, dat := range data {
		if dat.PopulationAmountFrom == 0 {
			continue
		}
		growthPct := float64(dat.PopulationAmountTo-dat.PopulationAmountFrom) / float64(dat.PopulationAmountFrom) * 100
		report = append(report, &domain.PopulationGrowthReport{
			OKTMO:               dat.OKTMO,
			YearFrom:            dat.YearFrom,
			YearTo:              dat.YearTo,
			PopulationGrowthPct: growthPct,
		})
	}
	return report
}
