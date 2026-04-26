// internal/handler/public/rosstat_test.go
package public_test

import (
	"backend/src/internal/config"
	"backend/src/internal/db/postgres"
	"backend/src/internal/dto"
	"backend/src/internal/handler/public"
	"backend/src/internal/model"
	repoImpl "backend/src/internal/repository/impl"
	serviceImpl "backend/src/internal/service/impl"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func cleanup(t *testing.T, conn *postgres.PostgresDBConnection) {
	db := conn.Get().(*gorm.DB)
	db.Exec("DELETE FROM midray.rosstat_age")
	db.Exec("DELETE FROM midray.rosstat")
	db.Exec("DELETE FROM midray.geo")
}

func TestGetRosstatHandler(t *testing.T) {
	godotenv.Load("../../../../.env")

	cfg := config.Load()
	conn := postgres.NewPostgresConnection(cfg.GetDBDSN())
	db := conn.Get().(*gorm.DB)
	cleanup(t, conn)

	geoModel := &model.Geo{
		Code:       45001,
		ParentCode: nil,
		Name:       "Челябинск",
		Level:      2,
	}
	if err := db.Create(geoModel).Error; err != nil {
		t.Fatalf("failed to create geo: %v", err)
	}

	birth2023 := 12000
	death2023 := 15000
	arrival2023 := 5000
	departure2023 := 3000
	male2023 := 560000
	female2023 := 640000
	landArea := 500.5
	avgSalary := 55000.0
	medical := 45
	schools := 120
	housing := 250.3

	rosstat := &model.Rosstat{
		Code:                45001,
		Year:                2023,
		PopulationAmount:    1200000,
		BirthAmount:         &birth2023,
		DeathAmount:         &death2023,
		ArrivalAmount:       &arrival2023,
		DepartureAmount:     &departure2023,
		MaleAmount:          &male2023,
		FemaleAmount:        &female2023,
		LandArea:            &landArea,
		AvgSalary:           &avgSalary,
		MedicalFacilities:   &medical,
		SchoolsCount:        &schools,
		HousingCommissioned: &housing,
	}
	if err := db.Create(rosstat).Error; err != nil {
		t.Fatalf("failed to create rosstat 2023: %v", err)
	}

	ages := []*model.RosstatAge{
		{RosstatID: rosstat.ID, Age: 0, MaleAmount: 5000, FemaleAmount: 4800},
		{RosstatID: rosstat.ID, Age: 25, MaleAmount: 7000, FemaleAmount: 6800},
		{RosstatID: rosstat.ID, Age: 65, MaleAmount: 4000, FemaleAmount: 6000},
	}
	for _, age := range ages {
		if err := db.Create(age).Error; err != nil {
			t.Fatalf("failed to create age: %v", err)
		}
	}

	rosstatRepo := repoImpl.NewRosstatRepository()
	rosstatAgeRepo := repoImpl.NewRosstatAgeRepository()
	rosstatService := serviceImpl.NewRosstatService(conn, rosstatRepo, rosstatAgeRepo)

	params := dto.RosstatParams{
		Codes: []int{45001},
		Fields: []string{"population", "birth",
			"death", "arrival", "departure", "male", "female", "land_area",
			"avg_salary", "medicial_facilities", "schools", "housing_commissioned", "age"},
	}

	response, err := public.GetRosstatHandler(nil, params, rosstatService)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response) == 0 {
		t.Fatal("expected non empty response")
	}
	if response[0].Code != 45001 {
		t.Fatalf("expected code 45001, got %d", response[0].Code)
	}
	if len(response[0].ByYear) != 1 {
		t.Fatalf("expected 1 year, got %d", len(response[0].ByYear))
	}

	year2023 := response[0].ByYear[0]
	if year2023.Population == nil || *year2023.Population != 1200000 {
		t.Fatalf("population: expected 1200000, got %v", year2023.Population)
	}
	if year2023.Birth == nil || *year2023.Birth != 12000 {
		t.Fatalf("birth: expected 12000, got %v", year2023.Birth)
	}
	if len(year2023.ByAge) != 3 {
		t.Fatalf("age: expected 3 records, got %d", len(year2023.ByAge))
	}
	if year2023.ByAge[0].Age != 0 || year2023.ByAge[0].Male != 5000 || year2023.ByAge[0].Female != 4800 {
		t.Fatalf("age0: expected 0/5000/4800, got %d/%d/%d", year2023.ByAge[0].Age, year2023.ByAge[0].Male, year2023.ByAge[0].Female)
	}

}
