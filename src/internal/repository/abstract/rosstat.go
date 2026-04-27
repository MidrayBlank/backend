package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IRosstatRepository interface {
	Upsert(conn abstract.IDBConnection, rosstat *domain.Rosstat) error
	UpsertBatch(conn abstract.IDBConnection, rosstats []*domain.Rosstat) error
	GetRosstatByCodes(conn abstract.IDBConnection, codes []int) ([]*domain.Rosstat, error)
}
