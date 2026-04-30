package domain

type RosstatAgeParsed struct {
	Age          int
	MaleAmount   int
	FemaleAmount int
}

type RosstatParsed struct {
	Code                int
	ParentCode          int
	Name                string
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
	ByAge               []*RosstatAgeParsed
}

type RosstatParsedSlice []*RosstatParsed
