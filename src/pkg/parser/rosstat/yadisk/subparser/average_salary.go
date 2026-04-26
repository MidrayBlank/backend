package subparser

import (
	"context"

	"backend/src/pkg/parser/rosstat/yadisk/dao"
	"backend/src/pkg/parser/rosstat/yadisk/extractor"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
)

type AverageSalarySubparser struct {
	BaseSubparser[dao.AverageSalaryExtracted]
}

func (p AverageSalarySubparser) Parse(ctx context.Context, storage *storage.Storage, url string) error {
	ext := extractor.NewAverageSalaryExtractor()
	return p.BaseParse(ctx, storage.SetAverageSalary, ext.Extract, url)
}
