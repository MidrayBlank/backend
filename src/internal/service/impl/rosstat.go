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

func (service *RosstatService) GetRosstatByCodes(codes []int) (domain.RosstatDataList, error) {
	rosstatInfos, err := service.rosstatRepo.GetRosstatByCodes(service.conn, codes)
	if err != nil {
		return nil, err
	}

	rosstatIDs := make([]int, len(rosstatInfos))
	for i, rosstatInfo := range rosstatInfos {
		rosstatIDs[i] = rosstatInfo.ID
	}

	rosstatAges, err := service.rosstatAgeRepo.GetRosstatAgeByRosstatIDs(service.conn, rosstatIDs)
	if err != nil {
		return nil, err
	}

	rosstatAgesData, err := service.toRosstatAgeData(rosstatAges)
	if err != nil {
		return nil, err
	}

	agesMap := make(map[int][]domain.RosstatAgeData)
	for _, rosstatAgeData := range rosstatAgesData {
		agesMap[rosstatAgeData.RosstatID] = append(agesMap[rosstatAgeData.RosstatID], rosstatAgeData)
	}

	return service.toRosstatDataList(rosstatInfos, agesMap)
}

func (service *RosstatService) toRosstatAgeData(rosstatAges []*domain.RosstatAge) ([]domain.RosstatAgeData, error) {
	result := make([]domain.RosstatAgeData, len(rosstatAges))

	for i, rosstatAge := range rosstatAges {
		if err := copier.Copy(&result[i], rosstatAge); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (service *RosstatService) toRosstatDataList(
	rosstatInfo []*domain.Rosstat,
	agesMap map[int][]domain.RosstatAgeData,
) (domain.RosstatDataList, error) {

	result := make([]domain.RosstatData, len(rosstatInfo))

	for i, info := range rosstatInfo {
		var data domain.RosstatData
		if err := copier.Copy(&data, info); err != nil {
			return nil, err
		}

		data.AgeData = agesMap[info.ID]
		result[i] = data
	}

	return result, nil
}
