package ai

import "github.com/gofiber/fiber/v3"

type IAIRouter interface {
	// SendRequestWithQuestion отправляет запрос по номеру вопроса (1,2,3)
	SendRequestWithQuestion(ctx fiber.Ctx, questionNumber int, inputData string) (string, error)

	// GetRequestsBalance возвращает количество оставшихся запросов
	GetRequestsBalance() int
}
