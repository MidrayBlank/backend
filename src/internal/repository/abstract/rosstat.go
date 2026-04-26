package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IRosstatRepository interface {
	Upsert(conn abstract.IDBConnection, rosstat *domain.Rosstat) error
	GetRosstatByCodes(conn abstract.IDBConnection, codes []string) ([]*domain.Rosstat, error)
}
