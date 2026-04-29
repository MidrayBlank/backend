package report

import "backend/src/pkg/report/async/subject_average/domain"

func buildSubjectRecord(subjectCode int, year int, municipalities domain.RosstatSlice) *domain.Rosstat {
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
		ByAge:               aggregateByAge(municipalities),
	}
}
