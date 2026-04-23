package abstract

import "context"

// Интерфейс по работе с счетчиком запросов
type IAIAPIRepository interface {

	// Создает новый счетчик(начинается с 0)
	CreateCounter(ctx context.Context, hash string) error

	// Увеличивает счеткик запросов
	IncrementRequestsCount(ctx context.Context, hash string) error

	// Получение кол-ва использованных запросов
	GetRequestsCount(ctx context.Context, hash string) (int, error)

	// Обнуляет счетчик запросов
	ResetRequestsCount(ctx context.Context, hash string) error
}
