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

func (modelObj *Geo) ToDomain() (*domain.Geo, error) {
	return ToDomain[Geo, domain.Geo](modelObj)
}

func (modelObj *Geo) ToModel(domainObj *domain.Geo) (*Geo, error) {
	return ToModel[Geo, domain.Geo](domainObj)
}

func (modelObj Geo) ToDomainSlice(modelObjs []Geo) ([]*domain.Geo, error) {
	return ToDomainSlice[Geo, domain.Geo](modelObjs)
}
