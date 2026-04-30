package domain

type PopulationDensityParams struct {
	OKTMO            int
	Year             int
	PopulationAmount int
	LandArea         float64
}

type PopulationDensityReport struct {
	OKTMO   int
	Year    int
	Density float64
}
