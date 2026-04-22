package impl

import (
	"backend/src/internal/domain"
)

type RosstatService struct {
}

func (service *RosstatService) GetRosstatByCodes(codes []int) domain.RosstatList {
	// TODO: connect to Repository and return real values
	// Mock just for tests
	result := make(domain.RosstatList, len(codes))

	for index := range result {
		result[index] = domain.RosstatGeo{Code: codes[index], Population: codes[index] * 2}
	}

	return result
}
