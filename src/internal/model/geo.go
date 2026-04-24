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
	if d == nil {
		return nil
	}
	return copier.Copy(m, d)
}

func (m Geo) ToDomainSlice(daos []Geo) ([]domain.Geo, error) {
	result := make([]domain.Geo, len(daos))
	for i, dao := range daos {
		domainItem, err := dao.ToDomain()
		if err != nil {
			return nil, err
		}
		result[i] = *domainItem
	}
	return result, nil
}
