package model

import (
	"backend/src/internal/domain"

	"github.com/jinzhu/copier"
)

type AiApi struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}

func (AiApi) TableName() string {
	return "ai_api"
}

func (m *AiApi) ToDomain() (*domain.AiApi, error) {
	var result domain.AiApi
	if err := copier.Copy(&result, m); err != nil {
		return nil, err
	}
	return &result, nil
}

func (m AiApi) ToDomainSlice(daos []AiApi) ([]domain.AiApi, error) {
	result := make([]domain.AiApi, len(daos))
	for i, dao := range daos {
		domainItem, err := dao.ToDomain()
		if err != nil {
			return nil, err
		}
		result[i] = *domainItem
	}
	return result, nil
}
