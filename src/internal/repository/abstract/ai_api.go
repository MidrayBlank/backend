package abstract

import "context"

// IAIAPIRepository - интерфейс по работе с счетчиком запросов
type IAIAPIRepository interface {

	// Upsert создает новый счетчик, если с таким hash его еще нет -> если есть, то счетчик увеличивается на 1
	Upsert(ctx context.Context, hash string) error

	// GetRequestsCount - получает количество использованных запросов
	GetRequestsCount(ctx context.Context, hash string) (int, error)

	// ResetRequestsCount - обнуляет счетчик запросов
	ResetRequestsCount(ctx context.Context, hash string) error
}
