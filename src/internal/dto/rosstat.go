package dto

type RosstatParams struct {
	Codes  []int    `query:"codes"`
	Fields []string `query:"fields" validate:"required,min=1"`
}

type RosstatByAge struct {
	Age    int `json:"age"`
	Male   int `json:"male"`
	Female int `json:"female"`
}

type RosstatByYear struct {
	Year                int            `json:"year"`
	Population          *int           `json:"population,omitempty"`
	Birth               *int           `json:"birth,omitempty"`
	Death               *int           `json:"death,omitempty"`
	Arrival             *int           `json:"arrival,omitempty"`
	Departure           *int           `json:"departure,omitempty"`
	Male                *int           `json:"male,omitempty"`
	Female              *int           `json:"female,omitempty"`
	LandArea            *int           `json:"land_area,omitempty"`
	AvgSalary           *float64       `json:"avg_salary,omitempty"`
	MedicalFacilities   *int           `json:"medical_facilities,omitempty"`
	SchoolsCount        *int           `json:"schools,omitempty"`
	HousingCommissioned *int           `json:"housing_commissioned,omitempty"`
	ByAge               []RosstatByAge `json:"by_age,omitempty"`
}

type Rosstat struct {
	Code   int             `json:"code"`
	ByYear []RosstatByYear `json:"by_year,omitempty"`
}

type RosstatResponse []Rosstat
