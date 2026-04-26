package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type ArrivalSubparser struct {
	BaseSubparser[dao.ArrivalExtracted]
}

func (p ArrivalSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewArrivalExtractor()
	return p.BaseParse(ctx, storage.SetArrival, ext.Extract, url)
}
