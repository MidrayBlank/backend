package impl

import (
	"backend/src/internal/model"
	"errors"

	"gorm.io/gorm"
)

type GeoRepository struct {
	db *gorm.DB
}

func NewGeoRepository(newDB *gorm.DB) *GeoRepository {
	return &GeoRepository{
		db: newDB,
	}
}

func (r *GeoRepository) GetGeoByCode(code int) (*model.Geo, error) {
	var result model.Geo

	err := r.db.Table("midray.geo").
		Where("code = ?", code).
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
