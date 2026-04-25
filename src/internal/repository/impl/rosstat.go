package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RosstatRepository struct{}

func NewRosstatRepository() *RosstatRepository {
	return &RosstatRepository{}
}

func (r *RosstatRepository) Upsert(conn abstract.IDBConnection, rosstat *domain.Rosstat) error {
	db := conn.Get().(*gorm.DB)

	rosstatDAO := &model.Rosstat{}
	rosstatDAO, err := rosstatDAO.ToModel(rosstat)

	if err != nil {
		return err
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "code"},
			{Name: "year"},
		},
		UpdateAll: true,
	}).Create(rosstatDAO).Error
}

func (r *RosstatRepository) GetRosstatByCodes(conn abstract.IDBConnection, codes []string) ([]*domain.Rosstat, error) {
	db := conn.Get().(*gorm.DB)

	var rosstatDAOs []model.Rosstat
	err := db.Where("code IN ?", codes).
		Find(&rosstatDAOs).Error
	if err != nil {
		return nil, err
	}

	var modelObj model.Rosstat
	return modelObj.ToDomainSlice(rosstatDAOs)
}
