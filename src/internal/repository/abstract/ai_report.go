package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type IAiReportRepository interface {
	Upsert(conn abstract.IDBConnection, code int, report string) error
	GetReportByCode(conn abstract.IDBConnection, code int) (*domain.AiReport, error)
}
