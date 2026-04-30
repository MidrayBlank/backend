package model

import "time"

type AsyncRequest struct {
	ID          int        `gorm:"primaryKey;autoIncrement"`
	Hash        string     `gorm:"column:hash;type:text;unique;not null"`
	Code        int        `gorm:"column:code;type:int;not null"`
	Status      int        `gorm:"column:status;type:int;not null"`
	RequestType int        `gorm:"column:request_type;type:int;not null;"`
	RequestData string     `gorm:"column:request_data;type:text"`
	Result      string     `gorm:"column:result;type:text"`
	Error       string     `gorm:"column:error;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`
	DeadlineAt  *time.Time `gorm:"column:deadline_at;type:timestamp;nullable"`
}
