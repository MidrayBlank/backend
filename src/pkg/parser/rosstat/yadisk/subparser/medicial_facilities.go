package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type MedicialFacilitiesSubparser struct {
	BaseSubparser[dao.MedicialFacilitiesExtracted]
}

func (p MedicialFacilitiesSubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewMedicialFacilitiesExtractor()
	return p.BaseParse(ctx, storage.SetMedicialFacilities, ext.Extract, url)
}
