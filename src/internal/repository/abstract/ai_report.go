package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IAIReportRepository interface {
	Put(conn abstract.IDBConnection, code int, report string) error
	Get(conn abstract.IDBConnection, code int) (*domain.AiReport, error)
	Delete(conn abstract.IDBConnection, code int) error
}
