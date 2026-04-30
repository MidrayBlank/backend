package impl

import (
	connection "backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"
	repository "backend/src/internal/repository/abstract"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
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
	request, err := service.asyncRequestRepo.GetSuccessRequestByCode(service.conn, code)
	if err != nil {
		return nil, err
	}

	if request == nil {
		return nil, nil
	}

	return request.ToDomain()
}

func (service *AiReportAsyncService) CreateReportRequest(code int, requestType int) (string, error) {
	exitingRequest, err := service.asyncRequestRepo.GetOneRequest(context.Background(), service.conn)
	if err != nil {
		return "", err
	}

	if exitingRequest != nil {
		return "", errors.New("request already exists")
	}

	hash := generateHash(code)

	request := model.NewAsyncRequest(hash, code, requestType)

	err = service.asyncRequestRepo.CreateRequest(service.conn, request)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (service *AiReportAsyncService) GetRequestStatusByHash(hash string) (*domain.AsyncRequest, error) {
	request, err := service.asyncRequestRepo.GetRequestByHash(service.conn, hash)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, nil
	}

	return request, nil
}

func generateHash(code int) string {
	hash := sha256.Sum256([]byte(strconv.Itoa(code)))
	return hex.EncodeToString(hash[:])
}
