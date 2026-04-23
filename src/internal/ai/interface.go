package ai

type IAIOpenRouter interface {
	// SendReques отправляет запрос в LLM по номеру вопроса
	SendRequest(requstData string) (string, error)
}
