package domain

type RosstatData struct {
	Code                int
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
	AgeData             []RosstatAgeData
}

type RosstatAgeData struct {
	Age          int
	MaleAmount   int
	FemaleAmount int
}

type RosstatDataList []RosstatData
