package ai

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type AIOpenRouter struct {
	client *OpenRouterClient
	model  string
}

// NewAIRouter создает объект AIRouter, создавая сначала объект OpenRouterClient
func NewAIOpenRouter(apiKey, model, baseURL string) *AIOpenRouter {
	return &AIOpenRouter{
		client: NewOpenRouterClient(apiKey, baseURL),
		model:  model,
	}
}

func (r *AIOpenRouter) UseApiKey(apiKey string) error {
	if strings.EqualFold(r.client.apiKey, apiKey) {
		return errors.New("этот API ключ уже стоит, нечего менять")
	}

	r.client.apiKey = apiKey
	return nil
}

// SendRequest создает запрос из полученной строки, затем обращается к OpenRouterClient для обработки и получения результата
func (r *AIOpenRouter) SendRequest(requstData string) (string, error) {
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
