package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
)

func GetRosstatHandler(ctx context.HandlerContext, params dto.RosstatParams, rosstatService service.IRosstatService) (dto.RosstatResponse, error) {
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

func buildRosstatResponse(codeMap map[int][]domain.RosstatByYear, fields []string) dto.RosstatResponse {
	result := make([]dto.Rosstat, 0, len(codeMap))
	for code, dataList := range codeMap {
		item := dto.Rosstat{
			Code:   code,
			ByYear: make([]dto.RosstatByYear, len(dataList)),
		}
		for i, data := range dataList {
			item.ByYear[i] = toRosstatByYear(data, fields)
		}

		result = append(result, item)
	}
	return result
}

func toRosstatByYear(data domain.RosstatByYear, fields []string) dto.RosstatByYear {
	byYear := dto.RosstatByYear{}

	for _, field := range fields {
		switch field {
		case "year":
			byYear.Year = data.Year

		case "population":
			byYear.Population = data.PopulationAmount

		case "birth":
			byYear.Birth = data.BirthAmount

		case "death":
			byYear.Death = data.DeathAmount

		case "arrival":
			byYear.Arrival = data.ArrivalAmount

		case "departure":
			byYear.Departure = data.DepartureAmount

		case "male":
			byYear.Male = data.MaleAmount

		case "female":
			byYear.Female = data.FemaleAmount

		case "land_area":
			byYear.LandArea = data.LandArea

		case "avg_salary":
			byYear.AvgSalary = data.AvgSalary

		case "medicial_facilities":
			byYear.MedicalFacilities = data.MedicalFacilities

		case "schools":
			byYear.SchoolsCount = data.SchoolsCount

		case "housing_commissioned":
			byYear.HousingCommissioned = data.HousingCommissioned

		case "age":
			byYear.ByAge = toRosstatByAgeSlice(data.AgeData)
		}
	}

	return byYear
}

func toRosstatByAge(age domain.RosstatByAge) dto.RosstatByAge {
	return dto.RosstatByAge{
		Age:    age.Age,
		Male:   age.MaleAmount,
		Female: age.FemaleAmount,
	}
}

func toRosstatByAgeSlice(ageData []domain.RosstatByAge) []dto.RosstatByAge {
	result := make([]dto.RosstatByAge, len(ageData))
	for i, age := range ageData {
		result[i] = toRosstatByAge(age)
	}
	return result
}
