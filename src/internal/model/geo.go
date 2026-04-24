package model

import (
	"backend/src/internal/domain"

	"github.com/jinzhu/copier"
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
	var result domain.Geo
	if err := copier.Copy(&result, m); err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *Geo) FromDomainToModel(d *domain.Geo) error {
	return copier.Copy(m, d)
}
