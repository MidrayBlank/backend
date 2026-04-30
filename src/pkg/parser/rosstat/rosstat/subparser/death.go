// subparser/death_subparser.go
package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/rosstat/config"
	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"backend/src/pkg/parser/rosstat/rosstat/extractor"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
)

type DeathSubparser struct {
	BaseSubparser[dao.DeathExtracted]
}

func NewDeathSubparser() *DeathSubparser {
	return &DeathSubparser{}
}

func (p *DeathSubparser) Parse(ctx context.Context, storage *storage.Storage) error {
	config := config.NewConfig()
	extractor := extractor.NewDeathExtractor(storage)

	return p.baseParse(ctx, storage, config.GetDeathIndicator, extractor.Extract)
}
