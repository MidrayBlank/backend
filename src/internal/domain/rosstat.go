package domain

type Rosstat struct {
	ID                  int
	Code                int
	Year                int
	PopulationAmount    *int
	BirthAmount         *int
	DeathAmount         *int
	ArrivalAmount       *int
	DepartureAmount     *int
	MaleAmount          *int
	FemaleAmount        *int
	LandArea            *int
	AvgSalary           *float64
	MedicalFacilities   *int
	SchoolsCount        *int
	HousingCommissioned *int
}

type RosstatByYear struct {
	Code                int
	Year                int
	PopulationAmount    *int
	BirthAmount         *int
	DeathAmount         *int
	ArrivalAmount       *int
	DepartureAmount     *int
	MaleAmount          *int
	FemaleAmount        *int
	LandArea            *int
	AvgSalary           *float64
	MedicalFacilities   *int
	SchoolsCount        *int
	HousingCommissioned *int
	AgeData             []RosstatByAge
}

type RosstatByAge struct {
	RosstatID    int
	Age          int
	MaleAmount   int
	FemaleAmount int
}

type RosstatList []RosstatByYear
