package impl

import (
	connection "backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	repository "backend/src/internal/repository/abstract"
)

type GeoService struct {
	conn    connection.IDBConnection
	geoRepo repository.IGeoRepository
}

func NewGeoService(conn connection.IDBConnection, geoRepo repository.IGeoRepository) *GeoService {
	return &GeoService{
		conn:    conn,
		geoRepo: geoRepo,
	}
}

func (service *GeoService) GetGeoAll() ([]domain.Geo, error) {
	geoList, err := service.geoRepo.GetGeoAll(service.conn)
	if err != nil {
		return nil, err
	}
	return toGeoSlice(geoList), nil
}

func toGeoSlice(geoList []*domain.Geo) []domain.Geo {
	result := make([]domain.Geo, len(geoList))
	for i, geo := range geoList {
		result[i] = *geo
	}
	return result
}
