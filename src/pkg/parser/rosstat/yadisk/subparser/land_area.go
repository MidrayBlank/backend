package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type LandAreaSubparser struct {
	BaseSubparser[dao.LandAreaExtracted]
}

func (p LandAreaSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewLandAreaExtractor()
	return p.BaseParse(ctx, storage.SetLandArea, ext.Extract, url)
}
