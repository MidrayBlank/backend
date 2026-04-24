package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AIAPIRepository struct{}

func NewAIAPIRepository() *AIAPIRepository {
	return &AIAPIRepository{}
}

func (r *AIAPIRepository) Upsert(conn abstract.IDBConnection, hash string) error {
	db := conn.Get().(*gorm.DB)

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hash"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"requests": gorm.Expr("requests + 1"),
		}),
	}).
		Create(&model.AIAPI{
			Hash:     hash,
			Requests: 0,
		}).Error
}

func (r *AIAPIRepository) GetRequestsCount(conn abstract.IDBConnection, hash string) (int, error) {
	var result model.AIAPI

	db := conn.Get().(*gorm.DB)
	err := db.Table("midray.ai_api").
		Where("hash = ?", hash).
		Select("requests").
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("there is no counter with this api key")
		}
		return 0, err
	}

	return result.Requests, nil
}

func (r *AIAPIRepository) ResetRequestsCount(conn abstract.IDBConnection, hash string) error {
	db := conn.Get().(*gorm.DB)
	result := db.Table("midray.ai_api").
		Where("hash = ?", hash).
		Update("requests", 0)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("there is no counter with this api key")
	}

	return nil
}
