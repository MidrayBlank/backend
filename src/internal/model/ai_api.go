package model

import (
	"backend/src/internal/domain"

	"github.com/jinzhu/copier"
)

type AIAPI struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}

func (AIAPI) TableName() string {
	return "ai_api"
}

func (m *AIAPI) ToDomain() (*domain.AIAPI, error) {
	var result domain.AIAPI
	if err := copier.Copy(&result, m); err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *AIAPI) FromDomainToModel(d *domain.AIAPI) error {
	if d == nil {
		return nil
	}
	return copier.Copy(m, d)
}

func (m AIAPI) ToDomainSlice(daos []AIAPI) ([]domain.AIAPI, error) {
	result := make([]domain.AIAPI, len(daos))
	for i, dao := range daos {
		domainItem, err := dao.ToDomain()
		if err != nil {
			return nil, err
		}
		result[i] = *domainItem
	}
	return result, nil
}
