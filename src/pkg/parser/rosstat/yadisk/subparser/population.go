package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type PopulationSubparser struct {
	BaseSubparser[dao.PopulationExtracted]
}

func (p PopulationSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewPopulationExtractor()
	return p.BaseParse(ctx, storage.SetPopulation, ext.Extract, url)
}
