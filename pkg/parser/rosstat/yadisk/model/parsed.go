package model

type RosstatParsed struct {
	ID                  int
	Code                string
	Year                int
	PopulationAmount    int
	BirthAmount         *int
	DeathAmount         *int
	ArrivalAmount       *int
	DepartureAmount     *int
	MaleAmount          *int
	FemaleAmount        *int
	LandArea            *float64
	AvgSalary           *float64
	MedicalFacilities   *int
	SchoolsCount        *int
	HousingCommissioned *float64
}

type RosstatAgeRecord struct {
	RosstatID    int
	Age          string
	MaleAmount   int
	FemaleAmount int
}

type RosstatParsedSlice []*RosstatParsed
