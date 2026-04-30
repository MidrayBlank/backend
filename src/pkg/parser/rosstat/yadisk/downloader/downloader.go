package downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func unzipFile(tempDir string, file *zip.File) (string, error) {
	destPath := filepath.Join(tempDir, file.Name)

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", err
	}

	rc, err := file.Open()
	if err != nil {
		return "", err
	}

	csvFile, err := os.Create(destPath)
	if err != nil {
		rc.Close()
		return "", err
	}

	_, err = io.Copy(csvFile, rc)
	csvFile.Close()
	rc.Close()

	if err != nil {
		return "", err
	}

	return destPath, nil
}

func DownloadAndUnzip(ctx context.Context, url string) ([]string, error) {
	fmt.Printf("Downloading: %s\n", url)

	tempDir, err := os.MkdirTemp(".", "rosstat_*")
	if err != nil {
		return nil, err
	}

	zipPath := filepath.Join(tempDir, "data.zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 60 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	out.Close()
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

	files := make([]*zip.File, 0)

	for _, file := range r.File {
		if strings.HasSuffix(file.Name, ".csv") && !strings.HasPrefix(filepath.Base(file.Name), "._") && !file.FileInfo().IsDir() {
			files = append(files, file)
		}
	}

	if len(files) == 1 {
		destPath, err := unzipFile(tempDir, files[0])

		if err != nil {
			return nil, err
		}

		csvPaths = append(csvPaths, destPath)
	} else {
		for _, file := range files {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			if strings.Contains(file.Name, "year") {
				destPath, err := unzipFile(tempDir, file)

				if err != nil {
					return nil, err
				}

				csvPaths = append(csvPaths, destPath)
			}
		}
	}

	if len(csvPaths) == 0 {
		return nil, fmt.Errorf("no CSV files found in zip")
	}

	fmt.Printf("Extracted %d CSV files\n", len(csvPaths))
	return csvPaths, nil
}

func CleanupTemp(path string) {
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
}
