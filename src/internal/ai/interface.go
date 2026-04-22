package ai

type IAIRouter interface {
	// SendReques отправляет запрос в LLM по номеру вопроса
	SendRequest(inputData string) (string, error)
}
