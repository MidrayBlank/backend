package abstract

import (
	"backend/src/internal/domain"
)

type IRawService interface {
	GetRosstatRawByCodes(codes []int) domain.RosstatList
}
