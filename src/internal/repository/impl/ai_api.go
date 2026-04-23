package impl

import (
	"backend/src/internal/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AIAPIRepository struct {
	db *gorm.DB
}

func NewAIAPIRepository(newDB *gorm.DB) *AIAPIRepository {
	return &AIAPIRepository{
		db: newDB,
	}
}

func (r *AIAPIRepository) Upsert(hash string) error {
	return r.db.Clauses(clause.OnConflict{
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

func (r *AIAPIRepository) GetRequestsCount(hash string) (int, error) {
	var result model.AIAPI

	err := r.db.Table("midray.ai_api").
		Where("hash = ?", hash).
		Select("requests").
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("счетчика с таким api не существует")
		}
		return 0, err
	}

	return result.Requests, nil
}

func (r *AIAPIRepository) ResetRequestsCount(hash string) error {
	result := r.db.Table("midray.ai_api").
		Where("hash = ?", hash).
		Update("requests", 0)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("счетчика с таким api не существует")
	}

	return nil
}
