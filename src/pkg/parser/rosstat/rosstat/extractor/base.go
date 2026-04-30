package extractor

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"backend/src/pkg/parser/rosstat/rosstat/state_manager"
)

type BaseCSVExtractor[T any] struct {
	file *os.File
}

func (r *BaseCSVExtractor[T]) getReader(filePath string) (*csv.Reader, error) {
	var err error
	r.file, err = os.Open(filePath)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(r.file)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'
	reader.LazyQuotes = true

	return reader, nil
}

func (r *BaseCSVExtractor[T]) closeReader() {
	r.file.Close()
}

func (r *BaseCSVExtractor[T]) ExtractRows(
	ctx context.Context,
	filePath string,
	manager *state_manager.StateManager[T],
) error {
	reader, err := r.getReader(filePath)

	if err != nil {
		return err
	}

	defer r.closeReader()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read error at file %s: %w", filePath, err)
		}

		manager.DoOnStateState(row)
	}
	return nil
}
