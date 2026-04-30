package domain

import "time"

type AsyncRequest struct {
	ID          int
	Hash        string
	Code        int
	Status      int
	RequestType int
	Result      string
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
