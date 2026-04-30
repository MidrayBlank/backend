package abstract

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/model"
	"context"
)

type IAsyncRequestRepository interface {
	GetOneRequest(ctx context.Context, conn abstract.IDBConnection) (*model.AsyncRequest, error)
	SetStatusById(ctx context.Context, conn abstract.IDBConnection, id int, status int) error
	SetStatusAndIncrementById(ctx context.Context, conn abstract.IDBConnection, id int, status int) error
	CloseTimeoutRequests(ctx context.Context, conn abstract.IDBConnection) error
}
