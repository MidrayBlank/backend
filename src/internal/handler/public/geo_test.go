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

func insertTestGeo(t *testing.T, db *gorm.DB) {
	federalSubjects := []model.Geo{
		{Code: 45000, ParentCode: nil, Name: "Челябинская область", Level: 1},
		{Code: 55000, ParentCode: nil, Name: "Свердловская область", Level: 1},
	}
	for _, federalSubject := range federalSubjects {
		if err := db.Create(&federalSubject).Error; err != nil {
			t.Fatalf("failed to create subject: %v", err)
		}
	}

	upperMunicipalities := []model.Geo{
		{Code: 45001, ParentCode: &[]int{45000}[0], Name: "Челябинск", Level: 2},
		{Code: 45002, ParentCode: &[]int{45000}[0], Name: "Магнитогорск", Level: 2},
		{Code: 55001, ParentCode: &[]int{55000}[0], Name: "Екатеринбург", Level: 2},
	}
	for _, upperMunicipality := range upperMunicipalities {
		if err := db.Create(&upperMunicipality).Error; err != nil {
			t.Fatalf("failed to create upper municipality: %v", err)
		}
	}

	lowerMunicipalities := []model.Geo{
		{Code: 45003, ParentCode: &[]int{45001}[0], Name: "Ленинский район", Level: 3},
		{Code: 45004, ParentCode: &[]int{45001}[0], Name: "Курчатовский район", Level: 3},
	}
	for _, lowerMunicipality := range lowerMunicipalities {
		if err := db.Create(&lowerMunicipality).Error; err != nil {
			t.Fatalf("failed to create lower municipality: %v", err)
		}
	}
}

func TestGetGeoHandler(t *testing.T) {
	godotenv.Load("../../../../.env")

	cfg := config.Load()
	conn := postgres.NewPostgresConnection(cfg.GetDBDSN())
	db := conn.Get().(*gorm.DB)
	cleanup(t, conn)

	insertTestGeo(t, db)

	geoRepo := repoImpl.NewGeoRepository()
	geoService := serviceImpl.NewGeoService(conn, geoRepo)

	response, err := public.GetGeoHandler(nil, geoService)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf("expected 2 federal subjects, got %d", len(response))
	}

	var sverdlovskRegion *dto.FederalSubject
	for i := range response {
		if response[i].Code == 55000 {
			sverdlovskRegion = &response[i]
			break
		}
	}

	if sverdlovskRegion == nil {
		t.Errorf("sverdlovsk region is not founded")
	}

	if sverdlovskRegion.Name != "Свердловская область" {
		t.Errorf("expected: Свердловская область, got: %s", sverdlovskRegion.Name)
	}

	if len(sverdlovskRegion.UpperMunicipalities) != 1 {
		t.Errorf("expected: 1(only Екатеринбург), got: %d", len(sverdlovskRegion.UpperMunicipalities))
	}

	if sverdlovskRegion.UpperMunicipalities[0].Name != "Екатеринбург" {
		t.Errorf("expected: Екатеринбург, got: %s", sverdlovskRegion.UpperMunicipalities[0].Name)
	}

	cleanup(t, conn)
}
