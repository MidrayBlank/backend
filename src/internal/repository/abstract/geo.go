package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IGeoRepository interface {
	Upsert(conn abstract.IDBConnection, geo *domain.Geo) error
	UpsertBatch(conn abstract.IDBConnection, geos []*domain.Geo) error
	GetGeoAll(conn abstract.IDBConnection) ([]*domain.Geo, error)
}
