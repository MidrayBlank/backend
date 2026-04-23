package impl

import (
	"backend/src/internal/model"
	"context"
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

// Upsert создает новый счетчик, если с таким hash его еще нет -> если есть, то счетчик увеличивается на 1
func (r *AIAPIRepository) Upsert(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
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

// Получение кол-ва использованных запросов
func (r *AIAPIRepository) GetRequestsCount(ctx context.Context, hash string) (int, error) {
	var result model.AIAPI

	err := r.db.WithContext(ctx).
		Table("midray.ai_api").
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

// Обнуляет счетчик запросов
func (r *AIAPIRepository) ResetRequestsCount(ctx context.Context, hash string) error {
	result := r.db.WithContext(ctx).
		Table("midray.ai_api").
		Where("hash = ?", hash).
		Update("requests", 0)

	if result.Error != nil {
		return result.Error
	}

	// Записи нет значит
	if result.RowsAffected == 0 {
		return errors.New("счетчика с таким api не существует")
	}

	return nil
}
