package abstract

import (
	"backend/pkg/parser/rosstat/yadisk/model"
	"context"
)

type ParseResult struct {
	Records []*model.DemographyRecord
}

type IRosstatParser interface {
	Parse(ctx context.Context) (*ParseResult, error)
}
