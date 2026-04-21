package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenRouterClient struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

// NewOpenRouterClient создает новый клиент OpenRouter
func NewOpenRouterClient(apiKey, baseURL string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second, // Максимальное время ответа от сервера
		},
		baseURL: baseURL,
	}
}

// ChatCompletion отправляет запрос к модели и получает ответ
func (c *OpenRouterClient) ChatCompletion(model string, messages []Message) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:       model,
		Messages:    messages,
		Stream:      false, // Получаем ответ одним соо, если true -> то по словам получаем ответ
		Temperature: 0.1,   // Уровень креативности(для нас лучше использовать мб >1.2)
		MaxTokens:   1000,
	}

	return c.sendRequest(reqBody)
}

// sendRequest внутренний метод для отправки запросов
func (c *OpenRouterClient) sendRequest(reqBody ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %w", err)
	}

	// Отправляем запрос; bytes.NewBuffer(jsonData) -> Позволяет работать с данными как с io.Reader или io.Writer
	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err = json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &chatResp, nil
}

// setHeaders устанавливает заголовки для запроса к OpenRouter
func (c *OpenRouterClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Опционально(для статистики на openrouter):
	req.Header.Set("HTTP-Referer", "http://localhost:8080") // Наш сайт
	req.Header.Set("X-Title", "Go OpenRouter Client")       // Название нашего сервиса

}
