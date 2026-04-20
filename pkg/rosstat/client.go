package rosstat

// Client представляет клиент для работы с API Росстата
type Client struct {
	config Config
}

// NewClient создаёт новый клиент для работы с API Росстата
func NewClient(cfg Config) *Client {
	return &Client{
		config: cfg,
	}
}

// DownloadCSV скачивает CSV с данными Росстата для указанных территорий и годов
func (c *Client) DownloadCSV(territoryCodes []int, yearFrom, yearTo int) ([]byte, error) {
	params := NewRequestParams(territoryCodes, yearFrom, yearTo)
	return Download(c.config, params)
}

// DownloadCSVWithParams скачивает CSV на основе готовых параметров запроса
func (c *Client) DownloadCSVWithParams(params RequestParams) ([]byte, error) {
	return Download(c.config, params)
}
