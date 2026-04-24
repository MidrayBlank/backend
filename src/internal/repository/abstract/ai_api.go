package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IAIAPIRepository interface {
	Insert(conn abstract.IDBConnection, token string) error
	GetAllRequestsCount(conn abstract.IDBConnection, tokens []string) ([]domain.AIAPI, error)
	IncreaseRequests(conn abstract.IDBConnection, token string) error
	ResetAllRequestsCount(conn abstract.IDBConnection) error
}
