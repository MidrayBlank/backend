package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type MaleFemaleAgeSubparser struct {
	BaseSubparser[dao.MaleFemaleAgeExtracted]
}

func (p MaleFemaleAgeSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewMaleFemaleAgeExtractor()
	return p.BaseParse(ctx, storage.SetMaleFemaleAge, ext.Extract, url)
}
