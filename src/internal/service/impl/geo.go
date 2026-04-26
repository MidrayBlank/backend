package impl

import (
	connection "backend/src/internal/db/abstract"
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
