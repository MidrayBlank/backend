package abstract

import "backend/src/internal/db/abstract"

type IAIAPIRepository interface {
	Upsert(conn abstract.IDBConnection, hash string) error

	GetRequestsCount(conn abstract.IDBConnection, hash string) (int, error)

	ResetRequestsCount(conn abstract.IDBConnection, hash string) error
}
