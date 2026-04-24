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

	dao := &model.Rosstat{
		ID:                  rosstat.ID,
		Code:                rosstat.Code,
		Year:                rosstat.Year,
		PopulationAmount:    rosstat.PopulationAmount,
		BirthAmount:         rosstat.BirthAmount,
		DeathAmount:         rosstat.DeathAmount,
		ArrivalAmount:       rosstat.ArrivalAmount,
		DepartureAmount:     rosstat.DepartureAmount,
		MaleAmount:          rosstat.MaleAmount,
		FemaleAmount:        rosstat.FemaleAmount,
		LandArea:            rosstat.LandArea,
		AvgSalary:           rosstat.AvgSalary,
		MedicalFacilities:   rosstat.MedicalFacilities,
		SchoolsCount:        rosstat.SchoolsCount,
		HousingCommissioned: rosstat.HousingCommissioned,
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "code"},
			{Name: "year"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"population_amout",
			"birth_amount",
			"death_amount",
			"arrival_amount",
			"departure_amount",
			"male_amount",
			"female_amount",
			"land_area",
			"avg_salary",
			"medical_facilities",
			"schools_count",
			"housing_commissioned",
		}),
	}).
		Create(dao).Error
}

func (r *RosstatRepository) GetRosstatByCodes(conn abstract.IDBConnection, codes []string) ([]domain.Rosstat, error) {
	db := conn.Get().(*gorm.DB)

	var rosstatDAOs []model.Rosstat
	err := db.Where("code IN ?", codes).
		Find(&rosstatDAOs).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Rosstat, len(rosstatDAOs))
	for i, dao := range rosstatDAOs {
		result[i] = domain.Rosstat{
			ID:                  dao.ID,
			Code:                dao.Code,
			Year:                dao.Year,
			PopulationAmount:    dao.PopulationAmount,
			BirthAmount:         dao.BirthAmount,
			DeathAmount:         dao.DeathAmount,
			ArrivalAmount:       dao.ArrivalAmount,
			DepartureAmount:     dao.DepartureAmount,
			MaleAmount:          dao.MaleAmount,
			FemaleAmount:        dao.FemaleAmount,
			LandArea:            dao.LandArea,
			AvgSalary:           dao.AvgSalary,
			MedicalFacilities:   dao.MedicalFacilities,
			SchoolsCount:        dao.SchoolsCount,
			HousingCommissioned: dao.HousingCommissioned,
		}
	}

	return result, nil
}
