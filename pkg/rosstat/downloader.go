package rosstat

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// Download отправляет запрос к API Росстата и возвращает CSV-данные
func Download(cfg Config, params RequestParams) ([]byte, error) {
	body := BuildRequestBody(params)
	// Создаём транспорт с отключенной проверкой сертификата
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Transport: tr,
	}
	req, err := http.NewRequest("POST", cfg.BaseURL, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Origin", "https://rosstat.gov.ru")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", cfg.Accept)
	req.Header.Set("Referer", cfg.Referer)
	req.Header.Set("User-Agent", cfg.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %d %s", resp.StatusCode, resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if !isValidCSVContentType(contentType) {
		return nil, fmt.Errorf("ожидался CSV, получен Content-Type: %s", contentType)
	}
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}
	if len(rawData) == 0 {
		return nil, fmt.Errorf("получен пустой ответ")
	}
	utf8Data, err := decodeWindows1251(rawData)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования: %w", err)
	}
	if !looksLikeCSV(utf8Data) {
		os.WriteFile("debug_response.txt", utf8Data, 0644)
		return nil, fmt.Errorf("ответ не похож на CSV (нет разделителей), сохранён в debug_response.txt")
	}
	return utf8Data, nil
}

// decodeWindows1251 преобразует данные из кодировки Windows-1251 в UTF-8
func decodeWindows1251(data []byte) ([]byte, error) {
	decoder := charmap.Windows1251.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err != nil {
		return nil, err
	}
	return utf8Data, nil
}

// isValidCSVContentType проверяет, что Content-Type соответствует CSV
func isValidCSVContentType(contentType string) bool {
	lower := strings.ToLower(contentType)
	allowed := []string{"text/csv", "application/csv", "application/x-csv", "text/plain"}
	for _, allowedType := range allowed {
		if strings.Contains(lower, allowedType) {
			return true
		}
	}
	return false
}

// looksLikeCSV проверяет, что данные похожи на CSV (имеют разделители)
func looksLikeCSV(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	checkLen := 500
	if len(data) < checkLen {
		checkLen = len(data)
	}
	sample := string(data[:checkLen])
	return strings.Contains(sample, ";") || strings.Contains(sample, ",")
}
