package abstract

import (
	"context"

	"backend/src/pkg/parser/rosstat/domain"
)

type IRosstatParser interface {
	Parse(ctx context.Context) ([]*domain.RosstatParsed, error)
}
