package impl

import (
	"backend/src/internal/async/request"
	"backend/src/internal/async/status"
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AsyncRequestRepository struct{}

func NewAsyncRequesRepository() *AsyncRequestRepository {
	return &AsyncRequestRepository{}
}

func (r *AsyncRequestRepository) GetOneRequest(ctx context.Context, conn abstract.IDBConnection) (*model.AsyncRequest, error) {
	db := conn.Get().(*gorm.DB)

	query := `
		SELECT * FROM midray.async_requests 
        WHERE status = ? 
          AND (deadline_at IS NULL OR deadline_at > ?)
        LIMIT 1 
        FOR UPDATE SKIP LOCKED
	`

	var request model.AsyncRequest
	err := db.WithContext(ctx).
		Raw(query, status.StatusQueued, time.Now()).
		Scan(&request).Error

	return &request, err
}

func (r *AsyncRequestRepository) SetStatusById(ctx context.Context, conn abstract.IDBConnection, id int, status int) error {
	db := conn.Get().(*gorm.DB)

	return db.WithContext(ctx).
		Model(&model.AsyncRequest{}).
		Where("request_type = ? AND id = ?", request.AIREPORT, id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *AsyncRequestRepository) SetStatusAndIncrementById(ctx context.Context,
	conn abstract.IDBConnection, id int, stat int) error {
	db := conn.Get().(*gorm.DB)

	return db.WithContext(ctx).
		Model(&model.AsyncRequest{}).
		Where("request_type = ? AND id = ? AND status NOT IN (?, ?)",
			request.AIREPORT, id, status.StatusSuccess, status.StatusFailed).
		Updates(map[string]interface{}{
			"status": gorm.Expr("CASE WHEN attempts + 1 >= 3 THEN ? ELSE ? END",
				status.StatusFailed, stat),
			"attempts":   gorm.Expr("attempts + 1"),
			"updated_at": time.Now(),
		}).Error

}

func (r *AsyncRequestRepository) CloseTimeoutRequests(ctx context.Context, conn abstract.IDBConnection) error {
	db := conn.Get().(*gorm.DB)

	return db.WithContext(ctx).
		Model(&model.AsyncRequest{}).
		Where("request_type = ? AND status = ? AND deadline_at IS NOT NULL AND deadline_at < ?",
			request.AIREPORT, status.StatusInProgress, time.Now()).
		Updates(map[string]interface{}{
			"status": gorm.Expr("CASE WHEN attempts >= 3 THEN ? ELSE ? END",
				status.StatusFailed, status.StatusQueued),
			"updated_at": time.Now(),
		}).Error
}

func (r *AsyncRequestRepository) CreateRequest(conn abstract.IDBConnection, request *model.AsyncRequest) error {
	db := conn.Get().(*gorm.DB)
	return db.Create(request).Error
}

func (r *AsyncRequestRepository) GetSuccessRequestByCode(conn abstract.IDBConnection, code int) (*model.AsyncRequest, error) {
	db := conn.Get().(*gorm.DB)

	var asyncRequest model.AsyncRequest
	err := db.
		Where("equest_type = ? AND code = ? AND status = ?", request.AIREPORT, code, status.StatusSuccess).
		First(&asyncRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &asyncRequest, nil
}

func (r *AsyncRequestRepository) GetRequestByHash(conn abstract.IDBConnection, hash string) (*domain.AsyncRequest, error) {
	db := conn.Get().(*gorm.DB)

	var request model.AsyncRequest
	err := db.Where("hash = ?", hash).First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return request.ToDomain()
}
