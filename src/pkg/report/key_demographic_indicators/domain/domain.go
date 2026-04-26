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
	BirthRate     float32
	DeathRate     float32
	NaturalRate   float32
	MigrationRate float32
}
