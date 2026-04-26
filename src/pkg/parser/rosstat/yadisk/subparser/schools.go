package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type SchoolsSubparser struct {
	BaseSubparser[dao.SchoolsExtracted]
}

func (p SchoolsSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewSchoolsExtractor()
	return p.BaseParse(ctx, storage.SetSchools, ext.Extract, url)
}
