package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
)

func GetGeoHandler(ctx context.HandlerContext, params dto.RosstatParams, rosstatService service.IRosstatService) (dto.RosstatResponse, error) {
	rosstatList, err := rosstatService.GetRosstatByCodes(params.Codes)
	if err != nil {
		return nil, err
	}

	codeMap := make(map[int][]domain.RosstatByYear)
	for _, data := range rosstatList {
		codeMap[data.Code] = append(codeMap[data.Code], data)
	}

	rosstatresponse := buildRosstatResponse(codeMap, params.Fields)

	return rosstatresponse, nil
}
