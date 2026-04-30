package report

import "backend/src/pkg/report/sync/key_demographic_indicators/domain"

type DemographicReportCompiler struct{}

func NewDemographicReportCompiler() *DemographicReportCompiler {
	return &DemographicReportCompiler{}
}

func (c *DemographicReportCompiler) Compile(data []*domain.DemographicParams) []*domain.DemographicReport {
	if len(data) == 0 {
		return []*domain.DemographicReport{}
	}
	report := make([]*domain.DemographicReport, 0, len(data))
	for _, dat := range data {
		if dat.PopulationAmount == 0 {
			report = append(report, &domain.DemographicReport{
				OKTMO:         dat.OKTMO,
				Year:          dat.Year,
				BirthRate:     0,
				DeathRate:     0,
				NaturalRate:   0,
				MigrationRate: 0,
			})
			continue
		}
		population := float64(dat.PopulationAmount)
		birthRate := float64(dat.BirthAmount) / population * 1000
		deathRate := float64(dat.DeathAmount) / population * 1000
		naturalRate := float64(dat.BirthAmount-dat.DeathAmount) / population * 1000
		migrationRate := float64(dat.ArrivalAmount-dat.DepartureAmount) / population * 1000
		report = append(report, &domain.DemographicReport{
			OKTMO:         dat.OKTMO,
			Year:          dat.Year,
			BirthRate:     birthRate,
			DeathRate:     deathRate,
			NaturalRate:   naturalRate,
			MigrationRate: migrationRate,
		})
	}
	return report
}
