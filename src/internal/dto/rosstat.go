package dto

type RosstatParams struct {
	Codes  []int    `query:"codes"`
	Fields []string `query:"fields"`
}

type RosstatGeo struct {
	Code       int `json:"code"`
	Population int `json:"population"`
}

type RosstatResponse []RosstatGeo
