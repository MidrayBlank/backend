package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IAiApiRepository interface {
	Insert(conn abstract.IDBConnection, token string) error
	GetAllRequestsCount(conn abstract.IDBConnection, tokens []string) ([]*domain.AiApi, error)
	IncreaseRequests(conn abstract.IDBConnection, token string) error
	ResetAllRequestsCount(conn abstract.IDBConnection) error
	InsertIfNotExist(conn abstract.IDBConnection, tokens []string) error
}
