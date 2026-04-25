package model

import (
	"backend/src/internal/domain"
)

type Rosstat struct {
	ID                  int      `gorm:"column:id;primaryKey;"`
	Code                string   `gorm:"column:code;type:int;not null;index"`
	Year                int      `gorm:"column:year;type:int;not null;index"`
	PopulationAmount    int      `gorm:"column:population_amout;type:int;not null"`
	BirthAmount         *int     `gorm:"column:birth_amount;type:int;nullable"`
	DeathAmount         *int     `gorm:"column:death_amount;type:int;nullable"`
	ArrivalAmount       *int     `gorm:"column:arrival_amount;type:int;nullable"`
	DepartureAmount     *int     `gorm:"column:departure_amount;type:int;nullable"`
	MaleAmount          *int     `gorm:"column:male_amount;type:int;nullable"`
	FemaleAmount        *int     `gorm:"column:female_amount;type:int;nullable"`
	LandArea            *float64 `gorm:"column:land_area;type:decimal(10,2);nullable"`
	AvgSalary           *float64 `gorm:"column:avg_salary;type:decimal(10,2);nullable"`
	MedicalFacilities   *int     `gorm:"column:medical_facilities;type:int;nullable"`
	SchoolsCount        *int     `gorm:"column:schools_count;type:int;nullable"`
	HousingCommissioned *float64 `gorm:"column:housing_commissioned;type:decimal(10,2);nullable"`
}

func (Rosstat) TableName() string {
	return "rosstat"
}

func (model *Rosstat) ToDomain() (*domain.Rosstat, error) {
	return ToDomain[Rosstat, domain.Rosstat](model)
}

func (model Rosstat) ToDomainSlice(models []Rosstat) ([]domain.Rosstat, error) {
	return ToDomainSlice[Rosstat, domain.Rosstat](models)
}
func (model *Rosstat) ToModel(d *domain.Rosstat) error {
	return ToModel[Rosstat, domain.Rosstat](model, d)
}
