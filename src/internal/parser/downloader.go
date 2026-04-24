package parser

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func DownloadAndUnzip(ctx context.Context, url string) (string, error) {
	fmt.Printf("Downloading: %s\n", url)

	tempDir, err := os.MkdirTemp("", "rosstat_*")
	if err != nil {
		return "", err
	}

	zipPath := filepath.Join(tempDir, "data.zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	total, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	reader := &progressReader{
		reader: resp.Body,
		total:  total,
		onUpdate: func(downloaded, total int64) {
			fmt.Printf("\rDownloaded: %.2f MB / %.2f MB", float64(downloaded)/1024/1024, float64(total)/1024/1024)
		},
	}

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", err
	}
	fmt.Println()

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var csvPath string
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".csv") && !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			csvPath = filepath.Join(tempDir, f.Name)
			csvFile, err := os.Create(csvPath)
			if err != nil {
				return "", err
			}
			defer csvFile.Close()

			_, err = io.Copy(csvFile, rc)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if csvPath == "" {
		return "", fmt.Errorf("no CSV file found in zip")
	}

	fmt.Printf("Extracted: %s\n", csvPath)
	return csvPath, nil
}

func DownloadAndUnzipAll(ctx context.Context, url string) ([]string, error) {
	fmt.Printf("Downloading: %s\n", url)

	tempDir, err := os.MkdirTemp("", "rosstat_*")
	if err != nil {
		return nil, err
	}

	zipPath := filepath.Join(tempDir, "data.zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// СОЗДАЁМ ЗАПРОС С КОНТЕКСТОМ
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 60 * time.Minute,
	}

	resp, err := client.Do(req) // используем req
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	total, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	reader := &progressReader{
		reader: resp.Body,
		total:  total,
		onUpdate: func(downloaded, total int64) {
			fmt.Printf("\rDownloaded: %.2f MB / %.2f MB", float64(downloaded)/1024/1024, float64(total)/1024/1024)
		},
	}

	_, err = io.Copy(out, reader)
	if err != nil {
		return nil, err
	}
	fmt.Println()

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var csvPaths []string

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".csv") && !f.FileInfo().IsDir() {
			destPath := filepath.Join(tempDir, f.Name)

			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return nil, err
			}

			rc, err := f.Open()
			if err != nil {
				return nil, err
			}

			csvFile, err := os.Create(destPath)
			if err != nil {
				rc.Close()
				return nil, err
			}

			_, err = io.Copy(csvFile, rc)
			csvFile.Close()
			rc.Close()

			if err != nil {
				return nil, err
			}

			csvPaths = append(csvPaths, destPath)
		}
	}

	if len(csvPaths) == 0 {
		return nil, fmt.Errorf("no CSV files found in zip")
	}

	fmt.Printf("Extracted %d CSV files\n", len(csvPaths))
	return csvPaths, nil
}

type progressReader struct {
	reader   io.Reader
	total    int64
	current  int64
	onUpdate func(downloaded, total int64)
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.current += int64(n)
	if p.onUpdate != nil {
		p.onUpdate(p.current, p.total)
	}
	return n, err
}

func CleanupTemp(path string) {
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
}

func CleanupTempAll(tempDir string) {
	if tempDir != "" {
		os.RemoveAll(tempDir)
	}
}

// NormalizeOktmo преобразует OKTMO к коду муниципального района (первые 5 цифр + 000)
func NormalizeOktmo(oktmo string) string {
	oktmo = strings.TrimSpace(oktmo)
	if len(oktmo) < 5 {
		return oktmo
	}
	return oktmo[:5] + "000"
}

// GetDistrictOktmo возвращает код района (первые 5 цифр + 000)
func GetDistrictOktmo(oktmo string) string {
	if len(oktmo) >= 5 {
		return oktmo[:5] + "000"
	}
	return oktmo
}

// AggregateByDistrict агрегирует по районам
func AggregateByDistrict(records *[]RosstatRawRecord) map[string]map[int]float64 {
	result := make(map[string]map[int]float64)
	for _, r := range *records {
		districtOktmo := GetDistrictOktmo(r.Oktmo)
		if result[districtOktmo] == nil {
			result[districtOktmo] = make(map[int]float64)
		}
		result[districtOktmo][r.Year] += r.Value
	}
	return result
}

// AggregateSimple агрегирует по точному OKTMO
func AggregateSimple(records *[]RosstatRawRecord) map[string]map[int]float64 {
	result := make(map[string]map[int]float64)
	for _, r := range *records {
		if result[r.Oktmo] == nil {
			result[r.Oktmo] = make(map[int]float64)
		}
		result[r.Oktmo][r.Year] += r.Value
	}
	return result
}

// AggregateAgeSexWithNormalize агрегирует половозрастные данные с нормализацией OKTMO
func AggregateAgeSexWithNormalize(records *[]AgeSexRawRecord) map[string]map[int]struct {
	Male   int
	Female int
} {
	result := make(map[string]map[int]struct {
		Male   int
		Female int
	})

	for _, r := range *records {
		normalizedOktmo := NormalizeOktmo(r.Oktmo)

		if result[normalizedOktmo] == nil {
			result[normalizedOktmo] = make(map[int]struct {
				Male   int
				Female int
			})
		}

		entry := result[normalizedOktmo][r.Year]
		if r.Sex == "Мужчины" {
			entry.Male += r.Value
		} else if r.Sex == "Женщины" {
			entry.Female += r.Value
		}
		result[normalizedOktmo][r.Year] = entry
	}

	return result
}

func IntPtr(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func Float64Ptr(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}
