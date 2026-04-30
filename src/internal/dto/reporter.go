package dto

type GrowthReportParams struct {
	Codes    string `query:"codes"`
	YearFrom int    `query:"year_from"`
	YearTo   int    `query:"year_to"`
}

type GrowthLeadersReportParams struct {
	Codes    string `query:"codes" validate:"required"`
	YearFrom int    `query:"year_from" validate:"required"`
	YearTo   int    `query:"year_to" validate:"required"`
	Limit    int    `query:"limit"`
}

type DensityReportParams struct {
	Codes string `query:"codes" validate:"required,min=1"`
	Year  int    `query:"year" validate:"required"`
}

type IndicatorsReportParams struct {
	Codes string `query:"codes" validate:"required,min=1"`
	Year  int    `query:"year" validate:"required"`
}

type GrowthReportResponse struct {
	Code                int     `json:"code"`
	YearFrom            int     `json:"year_from"`
	YearTo              int     `json:"year_to"`
	PopulationGrowthPct float64 `json:"population_growth_pct"`
}

type GrowthLeadersReportResponse struct {
	Code                int     `json:"code"`
	YearFrom            int     `json:"year_from"`
	YearTo              int     `json:"year_to"`
	PopulationGrowthPct float64 `json:"population_growth_pct"`
}

type GrowthLeadersResponse struct {
	TopGrowth  []GrowthLeadersReportResponse `json:"top_growth"`
	TopDecline []GrowthLeadersReportResponse `json:"top_decline"`
}

type DensityReportResponse struct {
	Code    int     `json:"code"`
	Density float64 `json:"density"`
}

type IndicatorsReportResponse struct {
	Code          int     `json:"code"`
	Year          int     `json:"year"`
	BirthRate     float64 `json:"birth_rate"`
	DeathRate     float64 `json:"death_rate"`
	NaturalRate   float64 `json:"natural_rate"`
	MigrationRate float64 `json:"migration_rate"`
}
