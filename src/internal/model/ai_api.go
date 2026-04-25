package model

import (
	"backend/src/internal/domain"
	"errors"
)

type AiApi struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}

func (AiApi) TableName() string {
	return "ai_api"
}

func (modelObj *AiApi) ToDomain(token string, requests int) (*domain.AiApi, error) {
	if token == "" {
		return nil, errors.New("token must not be empty")
	}

	if requests < 0 {
		return nil, errors.New("count of requsts must be positive")
	}

	return &domain.AiApi{
		Token:    token,
		Requests: requests,
	}, nil
}

func (modelObj *AiApi) ToDomainSlice(modelObjs []AiApi, tokensHashMap map[string]string) ([]domain.AiApi, error) {
	result := make([]domain.AiApi, len(modelObjs))
	for i, dao := range modelObjs {
		domainObj, err := dao.ToDomain(tokensHashMap[dao.Hash], dao.Requests)
		if err != nil {
			return nil, err
		}
		result[i] = *domainObj
	}

	return result, nil
}

func (modelObj *AiApi) ToModel(hash string) (*AiApi, error) {
	if hash == "" {
		return nil, errors.New("hash must not be empty")
	}

	return &AiApi{
		Hash:     hash,
		Requests: 0,
	}, nil
}
