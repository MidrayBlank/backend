package report

import "backend/src/pkg/report/async/subject_average/domain"

func buildRegionRecord(subjectCode int, year int, municipalities domain.RosstatSlice) *domain.Rosstat {

	allAgeGroups := make([][]*domain.RosstatAge, len(municipalities))

	for i, m := range municipalities {
		allAgeGroups[i] = m.ByAge
	}

	return &domain.Rosstat{
		Code:                subjectCode,
		SubjectCode:         subjectCode,
		Year:                year,
		Population:          aggregatePopulation(municipalities),
		Birth:               aggregateBirth(municipalities),
		Death:               aggregateDeath(municipalities),
		Arrival:             aggregateArrival(municipalities),
		Departure:           aggregateDeparture(municipalities),
		Male:                aggregateMale(municipalities),
		Female:              aggregateFemale(municipalities),
		LandArea:            aggregateLandArea(municipalities),
		MedicalFacilities:   aggregateMedicalFacilities(municipalities),
		Schools:             aggregateSchools(municipalities),
		HousingCommissioned: aggregateHousingCommissioned(municipalities),
		ByAge:               aggregateByAge(allAgeGroups),
	}
}
