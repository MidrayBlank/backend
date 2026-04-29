package report

import "backend/src/pkg/report/async/subject_average/domain"

func aggregateByAge(municipalities domain.RosstatSlice) []*domain.RosstatAge {
	ages := make(map[int]*domain.RosstatAge)

	for _, municip := range municipalities {
		if municip.ByAge == nil {
			continue
		}
		for _, ageItem := range municip.ByAge {
			if _, exists := ages[ageItem.Age]; !exists {
				ages[ageItem.Age] = &domain.RosstatAge{
					Age:    ageItem.Age,
					Male:   0,
					Female: 0,
				}
			}

			ages[ageItem.Age].Male += ageItem.Male
			ages[ageItem.Age].Female += ageItem.Female
		}
	}
	result := make([]*domain.RosstatAge, 0, len(ages))

	for _, val := range ages {
		result = append(result, val)
	}

	return result
}

func aggregatePopulation(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Population != nil {
			hasVal = true
			sum += *municip.Population
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateBirth(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Birth != nil {
			hasVal = true
			sum += *municip.Birth
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateDeath(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Death != nil {
			hasVal = true
			sum += *municip.Death
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateArrival(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Arrival != nil {
			hasVal = true
			sum += *municip.Arrival
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateDeparture(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Departure != nil {
			hasVal = true
			sum += *municip.Departure
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateMale(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Male != nil {
			hasVal = true
			sum += *municip.Male
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateFemale(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Female != nil {
			hasVal = true
			sum += *municip.Female
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateLandArea(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.LandArea != nil {
			hasVal = true
			sum += *municip.LandArea
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateMedicalFacilities(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.MedicalFacilities != nil {
			hasVal = true
			sum += *municip.MedicalFacilities
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateSchools(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.Schools != nil {
			hasVal = true
			sum += *municip.Schools
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateHousingCommissioned(municipalities domain.RosstatSlice) *int {
	hasVal := false
	sum := 0
	for _, municip := range municipalities {
		if municip.HousingCommissioned != nil {
			hasVal = true
			sum += *municip.HousingCommissioned
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}
