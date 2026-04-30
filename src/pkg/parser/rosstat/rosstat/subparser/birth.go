package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/rosstat/config"
	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"backend/src/pkg/parser/rosstat/rosstat/extractor"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
)

type BirthSubparser struct {
	BaseSubparser[dao.BirthExtracted]
}

func NewBirthSubparser() *BirthSubparser {
	return &BirthSubparser{}
}

func (p *BirthSubparser) Parse(ctx context.Context, storage *storage.Storage) error {
	config := config.NewConfig()
	extractor := extractor.NewBirthExtractor(storage)

	return p.baseParse(ctx, storage, config.GetBirthIndicator, extractor.Extract)
}
