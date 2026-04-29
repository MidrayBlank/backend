package domain

type RosstatAge struct {
	Age    int
	Male   int
	Female int
}

type Rosstat struct {
	Code                int
	SubjectCode         int
	Year                int
	Population          *int
	Birth               *int
	Death               *int
	Arrival             *int
	Departure           *int
	Male                *int
	Female              *int
	LandArea            *int
	AverageSalary       *float64
	MedicialFacilities  *int
	Schools             *int
	HousingCommissioned *int
	ByAge               []*RosstatAge
}

type RosstatSlice []*Rosstat
