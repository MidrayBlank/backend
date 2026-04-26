package report

import "backend/src/pkg/report/population_density/domain"

type PopulationDensityCalculator struct{}

func NewPopulationDensityCalculator() *PopulationDensityCalculator {
	return &PopulationDensityCalculator{}
}

func (c *PopulationDensityCalculator) CompileReport(data []*domain.PopulationDensityParams) []*domain.PopulationDensityReport {
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
