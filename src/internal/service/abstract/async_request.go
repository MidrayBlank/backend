package abstract

import (
	"backend/src/internal/domain"
)

type IAiReportAsyncService interface {
	GetReportByCode(code int) (*domain.AsyncRequest, error)
	CreateReportRequest(code int, requestType int) (string, error)
	GetRequestStatusByHash(hash string) (*domain.AsyncRequest, error)
}
