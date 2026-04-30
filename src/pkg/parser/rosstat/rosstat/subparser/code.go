package subparser

import (
	"context"
	"fmt"

	"backend/src/pkg/parser/rosstat/rosstat/config"
	"backend/src/pkg/parser/rosstat/rosstat/downloader"
	"backend/src/pkg/parser/rosstat/rosstat/extractor"
	"backend/src/pkg/parser/rosstat/rosstat/storage"

	"golang.org/x/sync/errgroup"
)

type CodeSubparser struct {
}

func NewCodeSubparser() *CodeSubparser {
	return &CodeSubparser{}
}

func (p *CodeSubparser) Parse(ctx context.Context, storage *storage.Storage) error {
	config := config.NewConfig()

	subjectsCount := len(config.SubjectCodes)

	startIndex := 0
	endIndex := config.DownloadHTMLBatchSize

	if endIndex >= subjectsCount {
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
				err := p.parseByCode(ctx, storage, subjectCode)

				if err != nil {
					return fmt.Errorf("Error occured by code %d: %w", subjectCode, err)
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

func (p *CodeSubparser) parseByCode(ctx context.Context, storage *storage.Storage, subjectCode int) error {
	fileName, err := downloader.DownloadHTML(ctx, subjectCode)

	if err != nil {
		return err
	}
	defer downloader.CleanupFile(fileName)

	extractor := extractor.NewCodeExtractor()
	codeExtracted, err := extractor.Extract(ctx, fileName)

	if err != nil {
		return err
	}

	storage.SetCode(subjectCode, codeExtracted)

	return nil
}
