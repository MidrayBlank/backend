package subparser

import (
	"context"
	"fmt"
	"path/filepath"

	"backend/src/pkg/parser/rosstat/yadisk/downloader"
)

type BaseSubparser[T any] struct {
}

func (_ *BaseSubparser[T]) BaseParse(
	ctx context.Context,
	storageSet func([]*T),
	extract func(context.Context, string) ([]*T, error),
	url string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	filePaths, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("arrival download error: %w", err)
	}

	var daos []*T

	for _, filePath := range filePaths {
		if err := ctx.Err(); err != nil {
			return err
		}

		records, err := extract(ctx, filePath)

		if err != nil {
			return fmt.Errorf("arrival extract error: %w", err)
		}

		daos = append(daos, records...)
	}

	if len(filePaths) > 1 {
		downloader.CleanupTemp(filepath.Dir(filePaths[0]))
	} else {
		downloader.CleanupTemp(filePaths[0])
	}
	storageSet(daos)

	return nil
}
