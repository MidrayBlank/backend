package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type DeathSubparser struct {
	BaseSubparser[dao.DeathExtracted]
}

func (p DeathSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewDeathExtractor()
	return p.BaseParse(ctx, storage.SetDeath, ext.Extract, url)
}
