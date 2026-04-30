package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/rosstat/config"
	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"backend/src/pkg/parser/rosstat/rosstat/extractor"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
)

type PopulationSubparser struct {
	BaseSubparser[dao.PopulationExtracted]
}

func NewPopulationSubparser() *PopulationSubparser {
	return &PopulationSubparser{}
}

func (p *PopulationSubparser) Parse(ctx context.Context, storage *storage.Storage) error {
	config := config.NewConfig()
	extractor := extractor.NewPopulationExtractor(storage)

	return p.baseParse(ctx, storage, config.GetPopulationIndicator, extractor.Extract)
}
