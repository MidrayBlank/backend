package subparser

import (
	"context"
	"fmt"

	"backend/src/pkg/parser/rosstat/rosstat/config"
	"backend/src/pkg/parser/rosstat/rosstat/downloader"
	"backend/src/pkg/parser/rosstat/rosstat/storage"

	"golang.org/x/sync/errgroup"
)

type BaseSubparser[T any] struct {
}

func (p *BaseSubparser[T]) baseParse(
	ctx context.Context,
	storage *storage.Storage,
	getIndicator func(int) int,
	extract func(context.Context, string) error,
) error {
	config := config.NewConfig()

	subjectsCount := len(config.SubjectCodes)

	startIndex := 0
	endIndex := config.DownloadHTMLBatchSize

	if endIndex > subjectsCount {
		endIndex = subjectsCount
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		g, ctx := errgroup.WithContext(ctx)

		for index := startIndex; index < endIndex; index++ {
			g.Go(func() error {
				subjectCode := config.SubjectCodes[index]
				err := p.parseByCode(
					ctx,
					storage,
					getIndicator,
					extract,
					subjectCode,
				)

				if err != nil {
					return fmt.Errorf(
						"Error occured by subject code %d, indicator %d: %w",
						subjectCode,
						getIndicator(subjectCode),
						err,
					)
				}

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		if endIndex >= subjectsCount {
			break
		}

		startIndex = endIndex
		endIndex += config.DownloadHTMLBatchSize

		if endIndex > subjectsCount {
			endIndex = subjectsCount
		}
	}

	return nil
}

func (p *BaseSubparser[T]) parseByCode(
	ctx context.Context,
	storage *storage.Storage,
	getIndicator func(int) int,
	extract func(context.Context, string) error,
	subjectCode int,
) error {
	fileName, err := downloader.DownloadCSV(
		ctx,
		subjectCode,
		getIndicator(subjectCode),
		storage.GetSubjectCodes(subjectCode),
	)

	if err != nil {
		return err
	}
	//defer downloader.CleanupFile(fileName)

	err = extract(ctx, fileName)

	if err != nil {
		return err
	}

	return nil
}
