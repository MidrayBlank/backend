package model

import "time"

type AsyncRequest struct {
	ID             int        `gorm:"primaryKey;autoIncrement"`
	Hash           string     `gorm:"column:hash;type:text;unique;not null"`
	Type           int        `gorm:"column:type;type:int;not null"`
	Parameter      int        `gorm:"column:parameter;type:int;not null"`
	Status         int        `gorm:"column:status;type:int;not null"`
	Attempts       int        `gorm:"column:attempts;type:int;not null"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
	DeadlineAt     *time.Time `gorm:"column:deadline_at;type:timestamp;nullable"`
	TimeoutSeconds int        `gorm:"column:timeout_seconds;type:int;not null"`
}
