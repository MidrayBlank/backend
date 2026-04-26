package extractor

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type RosstatCSVReader[T any] struct {
	file *os.File
}

func (r *RosstatCSVReader[T]) getReader(filePath string) (*csv.Reader, error) {
	var err error
	r.file, err = os.Open(filePath)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(r.file)
	reader.Comma = ';'

	_, err = reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty CSV file")
		}
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	return reader, nil
}

func (r *RosstatCSVReader[T]) closeReader() {
	r.file.Close()
}

func (r *RosstatCSVReader[T]) ExtractRows(
	ctx context.Context,
	filePath string,
	extractRow func([]string) *T,
) ([]*T, error) {
	reader, err := r.getReader(filePath)

	if err != nil {
		return nil, err
	}

	defer r.closeReader()

	var result []*T

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("read error at file %s: %w", filePath, err)
		}

		if extracted := extractRow(row); extracted != nil {
			result = append(result, extracted)
		}
	}
	return result, nil
}
