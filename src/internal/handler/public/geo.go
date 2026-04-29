package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
	"sort"
)

func GetGeoHandler(ctx context.HandlerContext, geoService service.IGeoService) (dto.GeoResponse, error) {
	geoList, err := geoService.GetGeoAll()
	if err != nil {
		return dto.GeoResponse{}, err
	}

	return buildGeoResponse(geoList)
}

func buildGeoResponse(geoList []domain.Geo) (dto.GeoResponse, error) {
	federalSubjects, childMap := groupGeoByParent(geoList)

	return dto.GeoResponse{
		FederalSubjects: buildFederalSubjects(federalSubjects, childMap),
	}, nil
}

func groupGeoByParent(geoList []domain.Geo) ([]domain.Geo, map[int][]domain.Geo) {
	var federalSubjects []domain.Geo
	childMap := make(map[int][]domain.Geo)

	for _, geo := range geoList {
		if geo.ParentCode == nil || *geo.ParentCode == 0 {
			federalSubjects = append(federalSubjects, geo)
		} else {
			childMap[*geo.ParentCode] = append(childMap[*geo.ParentCode], geo)
		}
	}
	return federalSubjects, childMap
}

func buildFederalSubjects(federalSubjects []domain.Geo, childMap map[int][]domain.Geo) []dto.FederalSubject {
	sortGeoByCode(federalSubjects)

	var result []dto.FederalSubject
	for _, federalSubject := range federalSubjects {
		federalSubjectDTO := dto.FederalSubject{
			Name:                federalSubject.Name,
			Code:                federalSubject.Code,
			UpperMunicipalities: buildUpperMunicipalities(childMap[federalSubject.Code], childMap),
		}

		result = append(result, federalSubjectDTO)
	}
	return result
}

func buildUpperMunicipalities(upperMunicipalities []domain.Geo, childMap map[int][]domain.Geo) []dto.UpperMunicipality {
	sortGeoByCode(upperMunicipalities)

	var result []dto.UpperMunicipality
	for _, upperMunicipality := range upperMunicipalities {
		upperMunicipalityDTO := dto.UpperMunicipality{
			Name:                upperMunicipality.Name,
			Code:                upperMunicipality.Code,
			LowerMunicipalities: buildLowerMunicipalities(childMap[upperMunicipality.Code]),
		}
		result = append(result, upperMunicipalityDTO)
	}
	return result
}

func buildLowerMunicipalities(lowerMunicipalities []domain.Geo) []dto.LowerMunicipality {
	sortGeoByCode(lowerMunicipalities)

	var result []dto.LowerMunicipality
	for _, lowerMunicipality := range lowerMunicipalities {
		lowerMunicipalityDTO := dto.LowerMunicipality{
			Name: lowerMunicipality.Name,
			Code: lowerMunicipality.Code,
		}
		result = append(result, lowerMunicipalityDTO)
	}
	return result
}

func sortGeoByCode(geoList []domain.Geo) {
	sort.Slice(geoList, func(i, j int) bool {
		return geoList[i].Code < geoList[j].Code
	})
}
