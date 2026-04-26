package impl

import (
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
)

type RosstatService struct {
	conn abstract.IDBConnection
}

func NewRosstatService(conn abstract.IDBConnection) *RosstatService {
	return &RosstatService{conn: conn}
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
