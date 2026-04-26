package domain

type PopulationDensityParams struct {
	OKTMO            int
	Year             int
	PopulationAmount int
	LandArea         float32
}

type PopulationDensityReport struct {
	OKTMO   int
	Year    int
	Density float32
}
