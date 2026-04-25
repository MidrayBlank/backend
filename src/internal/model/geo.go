package model

import (
	"backend/src/internal/domain"
)

type Geo struct {
	Code       int    `gorm:"column:code;primaryKey;type:int"`
	ParentCode *int   `gorm:"column:parent_code;type:int;nullable"`
	Name       string `gorm:"column:name;type:text;not null"`
	Level      int    `gorm:"column:level;type:int;not null"`
}

func (Geo) TableName() string {
	return "geo"
}

func (m *Geo) ToDomain() (*domain.Geo, error) {
	return ToDomain[Geo, domain.Geo](m)
}

func (m *Geo) ToModel(d *domain.Geo) error {
	return ToModel[Geo, domain.Geo](m, d)
}

func (m Geo) ToDomainSlice(daos []Geo) ([]domain.Geo, error) {
	return ToDomainSlice[Geo, domain.Geo](daos)
}
