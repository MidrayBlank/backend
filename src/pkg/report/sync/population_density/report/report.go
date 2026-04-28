package report

import "backend/src/pkg/report/sync/population_density/domain"

type PopulationDensityReportCompiler struct{}

func NewPopulationDensityReportCompiler() *PopulationDensityReportCompiler {
	return &PopulationDensityReportCompiler{}
}

func (c *PopulationDensityReportCompiler) Compile(data []*domain.PopulationDensityParams) []*domain.PopulationDensityReport {
	if len(data) == 0 {
		return []*domain.PopulationDensityReport{}
	}
	report := make([]*domain.PopulationDensityReport, 0, len(data))
	for _, dat := range data {
		var density float32
		if dat.LandArea > 0 {
			density = float32(dat.PopulationAmount) / dat.LandArea
		}
		report = append(report, &domain.PopulationDensityReport{
			OKTMO:   dat.OKTMO,
			Year:    dat.Year,
			Density: density,
		})
	}
	return report
}
