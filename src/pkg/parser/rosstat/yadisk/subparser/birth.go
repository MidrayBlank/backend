package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type BirthSubparser struct {
	BaseSubparser[dao.BirthExtracted]
}

func (p BirthSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewBirthExtractor()
	return p.BaseParse(ctx, storage.SetBirth, ext.Extract, url)
}
