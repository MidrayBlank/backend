package domain

type GrowthReport struct {
	Code                int
	YearFrom            int
	YearTo              int
	PopulationGrowthPct float64
}

type GrowthLeadersReport struct {
	Code                int
	YearFrom            int
	YearTo              int
	PopulationGrowthPct float64
}

type GrowthLeadersResponse struct {
	TopGrowth  []GrowthLeadersReport
	TopDecline []GrowthLeadersReport
}

type DensityReport struct {
	Code    int
	Year    int
	Density float64
}

type IndicatorsReport struct {
	Code          int
	Year          int
	BirthRate     float64
	DeathRate     float64
	NaturalRate   float64
	MigrationRate float64
}
