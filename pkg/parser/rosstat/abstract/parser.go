package abstract

import (
	"backend/pkg/parser/rosstat/yadisk/model"
	"context"
)

// ParseResult содержит результаты парсинга
type ParseResult struct {
	Records []*model.DemographyRecord
}

// IRosstatParser - единый интерфейс для парсера
type IRosstatParser interface {
	Parse(ctx context.Context) (*ParseResult, error)
}
