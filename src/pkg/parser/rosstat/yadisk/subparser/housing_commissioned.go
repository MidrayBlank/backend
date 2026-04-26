package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type HousingCommissionedSubparser struct {
	BaseSubparser[dao.HousingCommissionedExtracted]
}

func (p HousingCommissionedSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewHousingCommissionedExtractor()
	return p.BaseParse(ctx, storage.SetHousingCommissioned, ext.Extract, url)
}
