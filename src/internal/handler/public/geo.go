package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
)

func GetGeoHandler(ctx context.HandlerContext, geoService service.IGeoService) (dto.GeoResponse, error) {
	geoList, err := geoService.GetGeoByCodes(nil)
	if err != nil {
		return dto.GeoResponse{}, err
	}

	return buildGeoResponse(geoList)
}

func buildGeoResponse(geoList []*domain.Geo) (dto.GeoResponse, error) {

}
