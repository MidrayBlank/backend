package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IGeoRepository interface {
	Upsert(conn abstract.IDBConnection) error
	GetGeoByCodes(conn abstract.IDBConnection, codes []int) ([]*domain.Geo, error)
}
