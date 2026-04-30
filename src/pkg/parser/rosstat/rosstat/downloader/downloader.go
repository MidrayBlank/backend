package downloader

import (
	"backend/src/pkg/parser/rosstat/rosstat/config"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func DownloadCSV(ctx context.Context, subjectCode int, indicator int, codes []int) (string, error) {
	config := config.NewConfig()

	log.Printf("Downloading CSV: code=%d, indicator=%d\n", subjectCode, indicator)

	params := NewRequestParameters(
		[]int{indicator},
		getMunr(codes),
		codes,
		2026,
	)

	url := fmt.Sprintf("https://rosstat.gov.ru/dbscripts/munst/munst%02d/DBInet.cgi", subjectCode)

	req, err := http.NewRequestWithContext(ctx, "POST", url, params.buildRequestBody())
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", "https://rosstat.gov.ru/dbscripts/munst/munst87/DBInet.cgi")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36")

	client := &http.Client{
		Timeout: time.Duration(config.DownloadCSVTimeoutSeconds) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	attempts := 0
	var resp *http.Response

	for attempts < config.DownloadCSVMaxAttempts {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		attempts++
		time.Sleep(time.Duration(config.DownloadCSVTimeSleepSeconds) * time.Second)
	}

	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: code: %d, body: %s", resp.StatusCode, subjectCode, string(body))
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()
	utf8Body, err := decoder.Bytes(rawBody)
	if err != nil {
		return "", fmt.Errorf("decoding Windows-1251: %w", err)
	}

	fileName := getCSVFileName(subjectCode, indicator)

	if err := os.WriteFile(fileName, utf8Body, 0644); err != nil {
		return "", fmt.Errorf("error writing file %s: %w", fileName, err)
	}

	log.Printf("Saved CSV to %s\n", fileName)
	return fileName, nil
}

func DownloadHTML(ctx context.Context, subjectCode int) (string, error) {
	config := config.NewConfig()

	log.Printf("Downloading HTML: code=%d\n", subjectCode)

	urlStr := fmt.Sprintf("https://rosstat.gov.ru/dbscripts/munst/munst%02d/DBInet.cgi", subjectCode)

	body := url.Values{}
	body.Set("pl", strconv.Itoa(config.GetPopulationIndicator(subjectCode)))

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; OktmoGrabber/1.0)")

	client := &http.Client{
		Timeout: time.Duration(config.DownloadHTMLTimeoutSeconds) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	attempts := 0
	var resp *http.Response

	for attempts < config.DownloadHTMLMaxAttempts {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		attempts++
		time.Sleep(time.Duration(config.DownloadHTMLTimeSleepSeconds) * time.Second)
	}

	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: code: %d, body:  %s", resp.StatusCode, subjectCode, string(body))
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()
	utf8Body, err := decoder.Bytes(rawBody)
	if err != nil {
		return "", fmt.Errorf("decoding Windows-1251: %w", err)
	}

	fileName := getHTMLFileName(subjectCode)

	if err := os.WriteFile(fileName, utf8Body, 0644); err != nil {
		return "", fmt.Errorf("error writing file %s: %w", fileName, err)
	}

	log.Printf("Saved HTML to %s\n", fileName)
	return fileName, nil
}

func CleanupFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("Warning: failed to remove %s: %v", path, err)
	}
}

func getCSVFileName(code int, indicator int) string {
	now := time.Now().UnixNano()
	return filepath.Join(".", fmt.Sprintf("temp_%d_%d_%d.csv", code, indicator, now))
}

func getHTMLFileName(code int) string {
	now := time.Now().UnixNano()
	return filepath.Join(".", fmt.Sprintf("temp_%d_%d.html", code, now))
}

func getMunr(codes []int) []int {
	munr := make([]int, 0, len(codes))

	for _, code := range codes {
		if code%1000 == 0 {
			munr = append(munr, code)
		}
	}

	return munr
}
