package abstract

import (
	"context"

	"backend/pkg/parser/rosstat/yadisk/model"
)

type IRosstatParser interface {
	Parse(ctx context.Context) ([]*model.RosstatParsed, error)
}
