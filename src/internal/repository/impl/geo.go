package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GeoRepository struct{}

func NewGeoRepository() *GeoRepository {
	return &GeoRepository{}
}

func (r *GeoRepository) Upsert(conn abstract.IDBConnection, geo *domain.Geo) error {
	db := conn.Get().(*gorm.DB)

	geoDAO := &model.Geo{}
	geoDAO, err := geoDAO.ToModel(geo)

	if err != nil {
		return err
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		UpdateAll: true,
	}).
		Create(geoDAO).Error

}

func (r *GeoRepository) UpsertBatch(conn abstract.IDBConnection, geos []*domain.Geo) error {
	db := conn.Get().(*gorm.DB)

	geoDAO := &model.Geo{}
	geoDAOs := make([]*model.Geo, 0, len(geos))
	for _, geo := range geos {
		populated, err := geoDAO.ToModel(geo)
		if err != nil {
			return err
		}
		geoDAOs = append(geoDAOs, populated)
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		UpdateAll: true,
	}).
		Create(&geoDAOs).Error
}

func (r *GeoRepository) GetGeoAll(conn abstract.IDBConnection) ([]*domain.Geo, error) {
	db := conn.Get().(*gorm.DB)

	var geoDAOs []model.Geo
	err := db.Find(&geoDAOs).Error

	if err != nil {
		return nil, err
	}

	var modelObj model.Geo
	return modelObj.ToDomainSlice(geoDAOs)
}
