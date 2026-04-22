package ai

import (
	"fmt"
	"os"
)

type QuestionNumber int

type AIRouter struct {
	client *OpenRouterClient
	model  string
}

// NewAIRouter создает объект AIRouter, создавая сначала объект OpenRouterClient
func NewAIRouter(apiKey, model, baseURL string) *AIRouter {
	return &AIRouter{
		client: NewOpenRouterClient(apiKey, baseURL),
		model:  model,
	}
}

// SendRequest создает запрос из полученной строки, затем обращается к OpenRouterClient для обработки и получения результата
func (r *AIRouter) SendRequest(requstData string) (string, error) {
	// Создаем стуктуру запросов к LLM
	message := []Message{
		{
			Role:    "user",
			Content: requstData,
		},
	}

	// Создаем обращение к OpenRouter и получаем ответ
	response, err := r.client.ChatCompletion(r.model, message)
	if err != nil {
		// Проверка на timeout у OpenRouterClient
		if os.IsTimeout(err) {
			return "", fmt.Errorf("timeout, LLM не ответила за предоставленное ей время: %w", err)
		}
		return "", fmt.Errorf("ошибка обращения к LLM: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}
