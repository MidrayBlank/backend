package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IAiApiRepositoru interface {
	Insert(conn abstract.IDBConnection, token string) error
	GetAllRequestsCount(conn abstract.IDBConnection, tokens []string) ([]*domain.AiApi, error)
	IncreaseRequests(conn abstract.IDBConnection, token string) error
	ResetAllRequestsCount(conn abstract.IDBConnection) error
}
