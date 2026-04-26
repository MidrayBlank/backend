package impl

import (
	connection "backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	repository "backend/src/internal/repository/abstract"
	"errors"
)

type RosstatService struct {
	conn           connection.IDBConnection
	rosstatRepo    repository.IRosstatRepository
	rosstatAgeRepo repository.IRosstatAgeRepository
}

func NewRosstatService(
	conn connection.IDBConnection,
	rosstatRepo repository.IRosstatRepository,
	rosstatAgeRepository repository.IRosstatAgeRepository) *RosstatService {

	return &RosstatService{
		conn:           conn,
		rosstatRepo:    rosstatRepo,
		rosstatAgeRepo: rosstatAgeRepository,
	}
}

func (service *RosstatService) GetRosstatByCodes(codes []int) (domain.RosstatDataList, error) {
	if len(codes) == 0 {
		return nil, errors.New("no codes received")
	}

	// TODO: connect to Repository and return real values
	// Mock just for tests
	result := make(domain.RosstatDataList, len(codes))

	// for index := range result {
	// 	result[index] = domain.RosstatGeo{Code: codes[index], Population: codes[index] * 2}
	// }

	return result
}
