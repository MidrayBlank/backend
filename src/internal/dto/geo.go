package dto

type GeoResponse struct {
	Name            string           `json:"name"`
	Code            int              `json:"code"`
	FederalSubjects []FederalSubject `json:"federal_subjects"`
}

type FederalSubject struct {
	Name                string              `json:"name"`
	Code                int                 `json:"code"`
	UpperMunicipalities []UpperMunicipality `json:"upper_municipalities"`
}

type UpperMunicipality struct {
	Name                string              `json:"name"`
	Code                int                 `json:"code"`
	LowerMunicipalities []LowerMunicipality `json:"lower_municipalities"`
}

type LowerMunicipality struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}
