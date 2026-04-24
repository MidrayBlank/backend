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

	rosstatAgeDAO := model.RosstatAge{
		RosstatID:    rosstatAge.RosstatID,
		Age:          rosstatAge.Age,
		MaleAmount:   rosstatAge.MaleAmount,
		FemaleAmount: rosstatAge.FemaleAmount,
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "rosstat_id"},
			{Name: "age"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"male_amount",
			"female_amount",
		}),
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

	result := make([]domain.RosstatAge, len(rosstatAgeDAOs))
	for i, dao := range rosstatAgeDAOs {
		result[i] = domain.RosstatAge{
			RosstatID:    dao.RosstatID,
			Age:          dao.Age,
			MaleAmount:   dao.MaleAmount,
			FemaleAmount: dao.FemaleAmount,
		}
	}

	return result, nil
}
