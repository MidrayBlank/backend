package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IRosstatAgeRepository interface {
	Upsert(conn abstract.IDBConnection, rosstatAge *domain.RosstatAge) error
	GetRosstatAgeByCodes(conn abstract.IDBConnection, codes []int) ([]domain.RosstatAge, error)
}
