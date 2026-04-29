package abstract

import "backend/src/internal/domain"

type IAiReportAsyncService interface {
	GetReportByCode(code int) (*domain.AiReport, error)
	CreateReportRequest(code int, requestType int) (string, error)
	GetRequestStatusByHash(hash string) (*domain.AiReport, error)
}
