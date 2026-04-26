package model

import (
	"backend/src/internal/domain"
)

type RosstatAge struct {
	RosstatID    int `gorm:"column:rosstat_id;type:int;primaryKey"`
	Age          int `gorm:"column:age;type:int;primaryKey"`
	MaleAmount   int `gorm:"column:male_amount;type:int;not null"`
	FemaleAmount int `gorm:"column:female_amount;type:int;not null"`
}

func (modelObj *RosstatAge) ToDomain() (*domain.RosstatByAge, error) {
	return ToDomain[RosstatAge, domain.RosstatByAge](modelObj)
}

func (modelObj *RosstatAge) ToModel(domainObj *domain.RosstatByAge) (*RosstatAge, error) {
	return ToModel[RosstatAge, domain.RosstatByAge](domainObj)
}

func (modelObj RosstatAge) ToDomainSlice(modelObjs []RosstatAge) ([]*domain.RosstatByAge, error) {
	return ToDomainSlice[RosstatAge, domain.RosstatByAge](modelObjs)
}
