package domain

type PopulationGrowthParams struct {
	OKTMO                int
	YearFrom             int
	YearTo               int
	PopulationAmountFrom int
	PopulationAmountTo   int
}

type PopulationGrowthReport struct {
	OKTMO               int
	YearFrom            int
	YearTo              int
	PopulationGrowthPct float32
}
