package model

import (
	"backend/src/internal/domain"
)

type AiReport struct {
	Code   int    `gorm:"column:code;type:int;primaryKey"`
	Report string `gorm:"column:report;type:text;not null"`
}

func (modelObj *AiReport) ToDomain() (*domain.AiReport, error) {
	return ToDomain[AiReport, domain.AiReport](modelObj)
}

func (modelObj *AiReport) ToModel(domainObj *domain.AiReport) (*AiReport, error) {
	return ToModel[AiReport, domain.AiReport](domainObj)
}
