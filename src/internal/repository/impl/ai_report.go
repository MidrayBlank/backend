package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AiReportRepository struct{}

func NewAiReportRepository() *AiReportRepository {
	return &AiReportRepository{}
}

func (r *AiReportRepository) Upsert(conn abstract.IDBConnection, code int, report string) error {
	db := conn.Get().(*gorm.DB)

	aiReportDAO := model.NewAiReport(code, report)

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		UpdateAll: true,
	}).Create(aiReportDAO).Error
}

func (r *AiReportRepository) GetReportByCode(conn abstract.IDBConnection, code int) (*domain.AiReport, error) {
	db := conn.Get().(*gorm.DB)

	var aiReportDAO model.AiReport
	err := db.Where("code = ?", code).First(&aiReportDAO).Error

	if err != nil {
		return nil, err
	}

	return aiReportDAO.ToDomain()
}
