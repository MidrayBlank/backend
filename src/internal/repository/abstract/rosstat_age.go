package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IRosstatAgeRepository interface {
	Upsert(conn abstract.IDBConnection, rosstatAge *domain.RosstatAge) error
	GetRosstatAgeByRosstatIDs(conn abstract.IDBConnection, ids []int) ([]*domain.RosstatAge, error)
}
