package model

import (
	"backend/src/internal/async/status"
	"backend/src/internal/domain"
	"time"
)

type AsyncRequest struct {
	ID          int        `gorm:"primaryKey;autoIncrement"`
	Hash        string     `gorm:"column:hash;type:text;unique;not null"`
	Code        int        `gorm:"column:code;type:int;not null"`
	Status      int        `gorm:"column:status;type:int;not null"`
	RequestType int        `gorm:"column:request_type;type:int;not null;"`
	Result      string     `gorm:"column:result;type:text"`
	Attempts    int        `gorm:"column:attempts;type:int"`
	Error       string     `gorm:"column:error;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`
	DeadlineAt  *time.Time `gorm:"column:deadline_at;type:timestamp;nullable"`
}

func (modelObj *AsyncRequest) ToDomain() (*domain.AsyncRequest, error) {
	return ToDomain[AsyncRequest, domain.AsyncRequest](modelObj)
}

func NewAsyncRequest(hash string, code int, requestType int) *AsyncRequest {
	now := time.Now()
	return &AsyncRequest{
		Hash:        hash,
		Code:        code,
		Status:      status.StatusQueued,
		RequestType: requestType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
