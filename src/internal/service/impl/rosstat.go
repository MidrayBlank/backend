package impl

import (
	connection "backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	repository "backend/src/internal/repository/abstract"

	"github.com/jinzhu/copier"
)

type RosstatService struct {
	conn           connection.IDBConnection
	rosstatRepo    repository.IRosstatRepository
	rosstatAgeRepo repository.IRosstatAgeRepository
}

func NewRosstatService(
	conn connection.IDBConnection,
	rosstatRepo repository.IRosstatRepository,
	rosstatAgeRepository repository.IRosstatAgeRepository,
) *RosstatService {

	return &RosstatService{
		conn:           conn,
		rosstatRepo:    rosstatRepo,
		rosstatAgeRepo: rosstatAgeRepository,
	}
}

func (service *RosstatService) GetRosstatByCodes(codes []int) (domain.RosstatList, error) {
	rosstatByYears, err := service.rosstatRepo.GetRosstatByCodes(service.conn, codes)
	if err != nil {
		return nil, err
	}

	rosstatIDs := make([]int, len(rosstatByYears))
	for i, rosstatInfo := range rosstatByYears {
		rosstatIDs[i] = rosstatInfo.ID
	}

	rosstatAges, err := service.rosstatAgeRepo.GetRosstatAgeByRosstatIDs(service.conn, rosstatIDs)
	if err != nil {
		return nil, err
	}

	agesMap := make(map[int][]*domain.RosstatByAge)
	for _, rosstatAge := range rosstatAges {
		agesMap[rosstatAge.RosstatID] = append(agesMap[rosstatAge.RosstatID], rosstatAge)
	}

	return service.toRosstatList(rosstatByYears, agesMap)
}

func (service *RosstatService) toRosstatList(
	rosstatByYears []*domain.Rosstat,
	agesMap map[int][]*domain.RosstatByAge,
) (domain.RosstatList, error) {

	result := make([]domain.RosstatByYear, len(rosstatByYears))

	for i, info := range rosstatByYears {
		var data domain.RosstatByYear
		if err := copier.Copy(&data, info); err != nil {
			return nil, err
		}

		ageDataList := make([]domain.RosstatByAge, len(agesMap[info.ID]))
		for j, age := range agesMap[info.ID] {
			ageDataList[j] = *age
		}

		data.AgeData = ageDataList
		result[i] = data
	}

	return result, nil
}
