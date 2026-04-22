package dto

type RosstatParams struct {
	Codes  []int    `query:"codes"`
	Fields []string `query:"fields" validate:"required,min=1"`
}

type RosstatGeo struct {
	Code       int `json:"code"`
	Population int `json:"population"`
}

type RosstatResponse []RosstatGeo
