package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RosstatAgeRepository struct{}

func NewRosstatAgeRepository() *RosstatAgeRepository {
	return &RosstatAgeRepository{}
}

func (r *RosstatAgeRepository) Upsert(conn abstract.IDBConnection, rosstatAge *domain.RosstatAge) error {
	db := conn.Get().(*gorm.DB)

	rosstatAgeDAO := &model.RosstatAge{}
	if err := rosstatAgeDAO.FromDomainToModel(rosstatAge); err != nil {
		return err
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "rosstat_id"},
			{Name: "age"},
		},
		UpdateAll: true,
	}).
		Create(rosstatAgeDAO).Error
}

func (r *RosstatAgeRepository) GetRosstatAgeByRosstatIDs(conn abstract.IDBConnection, ids []int) ([]domain.RosstatAge, error) {
	db := conn.Get().(*gorm.DB)

	var rosstatAgeDAOs []model.RosstatAge
	err := db.Where("rosstat_id IN ?", ids).
		Find(&rosstatAgeDAOs).Error

	if err != nil {
		return nil, err
	}

	var m model.RosstatAge
	return m.ToDomainSlice(rosstatAgeDAOs)
}
