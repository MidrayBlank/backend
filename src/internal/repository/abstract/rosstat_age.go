package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IRosstatAgeRepository interface {
	Upsert(conn abstract.IDBConnection, rosstatAge *domain.RosstatByAge) error
	GetRosstatAgeByRosstatIDs(conn abstract.IDBConnection, ids []int) ([]*domain.RosstatByAge, error)
}
