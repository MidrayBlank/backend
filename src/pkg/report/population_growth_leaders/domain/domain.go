package domain

type GrowthLeadersParams struct {
	OKTMO                int
	YearFrom             int
	YearTo               int
	PopulationAmountFrom int
	PopulationAmountTo   int
	Limit                int
}

type GrowthLeadersReport struct {
	OKTMO               int
	YearFrom            int
	YearTo              int
	PopulationGrowthPct float64
}
