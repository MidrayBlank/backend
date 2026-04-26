package abstract

import (
	"backend/src/internal/domain"
)

type IGeoService interface {
	GetGeoByCodes(codes []int) ([]*domain.Geo, error)
}
