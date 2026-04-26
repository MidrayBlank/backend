package abstract

import (
	"backend/src/internal/domain"
)

type IGeoService interface {
	GetGeoAll() ([]domain.Geo, error)
}
