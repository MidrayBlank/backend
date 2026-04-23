package abstract

import (
	"backend/src/internal/model"
)

type IGeoRepository interface {
	GetGeoByCodes(code int) ([]model.Geo, error)

	UpsertGeo(geo *model.Geo) error
}
