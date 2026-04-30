package impl

import (
	connection "backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	repository "backend/src/internal/repository/abstract"
)

type AiReportAsyncService struct {
	conn             connection.IDBConnection
	asyncRequestRepo repository.IAsyncRequestRepository
}

func NewAiReportAsyncService(conn connection.IDBConnection,
	asyncRequestRepo repository.IAsyncRequestRepository,
) *AiReportAsyncService {
	return &AiReportAsyncService{
		conn:             conn,
		asyncRequestRepo: asyncRequestRepo,
	}
}

func (service *AiReportAsyncService) GetReportByCode(code int) (*domain.AsyncRequest, error) {

}

func (service *AiReportAsyncService) CreateReportRequest(code int, requestType int) (string, error) {

}

func (service *AiReportAsyncService) GetRequestStatusByHash(hash string) (*domain.AsyncRequest, error) {

}
