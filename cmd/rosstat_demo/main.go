package main

import (
	"fmt"
	"os"
	"rosstat-parser/pkg/rosstat"
)

func main() {
	// Конфигурация клиента
	cfg := rosstat.Config{
		BaseURL:   "https://rosstat.gov.ru/dbscripts/munst/munst87/DBInet.cgi",
		Timeout:   60,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		Referer:   "https://rosstat.gov.ru/dbscripts/munst/munst87/DBInet.cgi",
		Accept:    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}

	// Создаём клиента
	client := rosstat.NewClient(cfg)

	// Коды муниципальных образований (ОКТМО)
	codes := []int{
		87602000, 87604000, 87608000, 87612000, 87616000, 87620000,
		87624000, 87626000, 87628000, 87632000, 87636000, 87640000,
		87644000, 87648000, 87652000, 87700000, 87600000, 87500000,
	}

	// Скачиваем данные
	fmt.Println("Скачивание CSV...")
	data, err := client.DownloadCSV(codes, 2009, 2025)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	// Сохраняем в файл
	if err := os.WriteFile("output.csv", data, 0644); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Успешно! Сохранено %d байт в output.csv\n", len(data))
}
