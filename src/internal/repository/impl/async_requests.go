package impl

import (
	"backend/src/internal/async/status"
	"backend/src/internal/db/abstract"
	"backend/src/internal/model"
	"context"
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
		Where("status = ? AND id = ?", id).
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
		Where("id = ? AND status NOT IN (?, ?)",
			id, status.StatusSuccess, status.StatusFailed).
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
		Where("status = ? AND deadline_at IS NOT NULL AND deadline_at < ?",
			status.StatusInProgress, time.Now()).
		Updates(map[string]interface{}{
			"status": gorm.Expr("CASE WHEN attempts >= 3 THEN ? ELSE ? END",
				status.StatusFailed, status.StatusQueued),
			"updated_at": time.Now(),
		}).Error
}

func (r *AsyncRequestRepository) CreateRequest(ctx context.Context, conn abstract.IDBConnection) error {
	db := conn.Get().(*gorm.DB)

}
