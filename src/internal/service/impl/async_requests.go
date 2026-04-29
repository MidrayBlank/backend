package impl

import (
	connection "backend/src/internal/db/abstract"
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
