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

func (RosstatAge) TableName() string {
	return "rosstat_age"
}

func (modelObj *RosstatAge) ToDomain() (*domain.RosstatAge, error) {
	return ToDomain[RosstatAge, domain.RosstatAge](modelObj)
}

func (modelObj *RosstatAge) ToModel(domainObj *domain.RosstatAge) (*RosstatAge, error) {
	return ToModel[RosstatAge, domain.RosstatAge](domainObj)
}

func (modelObj RosstatAge) ToDomainSlice(modelObjs []RosstatAge) ([]*domain.RosstatAge, error) {
	return ToDomainSlice[RosstatAge, domain.RosstatAge](modelObjs)
}
