package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
)

func GetRosstatInfo(ctx context.HandlerContext, params dto.RosstatParams, rosstatService service.IRosstatService) (dto.RosstatResponse, error) {
	rosstatList := rosstatService.GetRosstatByCodes(params.Codes)

	for _, field := range params.Fields {
		if field == "population" {
			result := make([]dto.RosstatGeo, len(rosstatList))

			for index, value := range rosstatList {
				result[index] = dto.RosstatGeo{Code: value.Code, Population: value.Population}
			}

			return result, nil
		}
	}

	return dto.RosstatResponse{}, nil
}
