package domain

type DemographicParams struct {
	OKTMO            int
	Year             int
	PopulationAmount int
	BirthAmount      int
	DeathAmount      int
	ArrivalAmount    int
	DepartureAmount  int
}

type DemographicReport struct {
	OKTMO         int
	Year          int
	BirthRate     float64
	DeathRate     float64
	NaturalRate   float64
	MigrationRate float64
}
