package dao

type PopulationExtracted struct {
	Code       int
	Year       int
	Population int
	RuralUsed  bool
	UrbanUsed  bool
}

type BirthExtracted struct {
	Code  int
	Year  int
	Birth int
}

type DeathExtracted struct {
	Code  int
	Year  int
	Death int
}

type ArrivalExtracted struct {
	Code    int
	Year    int
	Arrival int
}

type DepartureExtracted struct {
	Code      int
	Year      int
	Departure int
}

type LandAreaExtracted struct {
	Code     int
	Year     int
	LandArea int
}

type AverageSalaryExtracted struct {
	Code          int
	Year          int
	AverageSalary float64
}

type MedicialFacilitiesExtracted struct {
	Code               int
	Year               int
	MedicialFacilities int
}

type SchoolsExtracted struct {
	Code    int
	Year    int
	Schools int
}

type HousingCommissionedExtracted struct {
	Code                int
	Year                int
	HousingCommissioned int
}

type MaleFemaleAgeExtracted struct {
	Code         int
	Year         int
	Age          int
	MaleAmount   int
	FemaleAmount int
}
