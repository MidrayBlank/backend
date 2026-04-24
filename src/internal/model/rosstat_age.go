package model

import (
	"backend/src/internal/domain"

	"github.com/jinzhu/copier"
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

func (m *RosstatAge) ToDomain() (*domain.RosstatAge, error) {
	var result domain.RosstatAge
	if err := copier.Copy(&result, m); err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *RosstatAge) FromDomainToModel(d *domain.RosstatAge) error {
	return copier.Copy(m, d)
}

func (m RosstatAge) ToDomainSlice(daos []RosstatAge) ([]domain.RosstatAge, error) {
	result := make([]domain.RosstatAge, len(daos))
	for i, dao := range daos {
		domainItem, err := dao.ToDomain()
		if err != nil {
			return nil, err
		}
		result[i] = *domainItem
	}
	return result, nil
}
