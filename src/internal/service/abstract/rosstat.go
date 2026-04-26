package abstract

import (
	"backend/src/internal/domain"
)

type IRosstatService interface {
	GetRosstatByCodes(codes []int) (domain.RosstatDataList, error)
}
