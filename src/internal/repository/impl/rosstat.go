package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"time"

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

func (r *RosstatRepository) UpsertBatch(conn abstract.IDBConnection, rosstats []*domain.Rosstat) error {
	db := conn.Get().(*gorm.DB)

	rosstatDAO := &model.Rosstat{}
	rosstatDAOs := make([]*model.Rosstat, 0, len(rosstats))
	for _, rs := range rosstats {
		populated, err := rosstatDAO.ToModel(rs)
		if err != nil {
			return err
		}
		rosstatDAOs = append(rosstatDAOs, populated)
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "code"},
			{Name: "year"},
		},
		UpdateAll: true,
	}).
		Create(&rosstatDAOs).Error
}

func (r *RosstatRepository) GetRosstatByCodes(conn abstract.IDBConnection, codes []int) ([]*domain.Rosstat, error) {
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

func (r *RosstatRepository) GetRosstatByCodeForLastFiveYears(
	conn abstract.IDBConnection,
	code int,
) ([]*domain.Rosstat, error) {
	db := conn.Get().(*gorm.DB)

	currentYear := time.Now().Year()
	startYear := currentYear - 5

	var rosstatDAOs []model.Rosstat
	err := db.Where("code = ? AND year >= ?", code, startYear).
		Find(&rosstatDAOs).Error

	if err != nil {
		return nil, err
	}

	var modelObj model.Rosstat
	return modelObj.ToDomainSlice(rosstatDAOs)
}
