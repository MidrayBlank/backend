package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type DepartureSubparser struct {
	BaseSubparser[dao.DepartureExtracted]
}

func (p DepartureSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewDepartureExtractor()
	return p.BaseParse(ctx, storage.SetDeparture, ext.Extract, url)
}
