package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
)

func RosstatRawHandler(ctx context.HandlerContext, params dto.RosstatParams, rawService service.IRawService) (dto.RosstatResponse, error) {
	populationList := rawService.GetRosstatRawByCodes(params.Codes)

	for _, field := range params.Fields {
		if field == "population" {
			result := make([]dto.RosstatGeo, len(populationList))

			for index, value := range populationList {
				result[index] = dto.RosstatGeo{Code: value.Code, Population: value.Population}
			}

			return result, nil
		}
	}

	return dto.RosstatResponse{}, nil
}
