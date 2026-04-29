package aggregator

import (
	"backend/src/pkg/report/async/domain"
	"context"
)

type RegionAggregator struct{}

func NewRegionAggregator() *RegionAggregator {
	return &RegionAggregator{}
}

func (r *RegionAggregator) Aggregate(ctx context.Context, data domain.RosstatSlice) (domain.RosstatSlice, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if len(data) == 0 {
		return domain.RosstatSlice{}, nil
	}
	groups, err := groupBySubjectCode(ctx, data)
	if err != nil {
		return nil, err
	}
	result := make(domain.RosstatSlice, 0, len(groups))
	for subcode, municip := range groups {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		year := 0
		if len(municip) > 0 {
			year = municip[0].Year
		}
		regionRecord := buildRegionRecord(subcode, year, municip)
		result = append(result, regionRecord)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return result, nil
}

func aggregateFloat(values []*float64) *float64 {
	hasVal := false
	var sum float64 = 0
	for _, val := range values {
		if val != nil {
			hasVal = true
			sum += *val
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func aggregateInt(values []*int) *int {
	hasVal := false
	sum := 0
	for _, val := range values {
		if val != nil {
			hasVal = true
			sum += *val
		}
	}
	if !hasVal {
		return nil
	}
	return &sum
}

func groupBySubjectCode(ctx context.Context, data domain.RosstatSlice) (map[int]domain.RosstatSlice, error) {
	result := make(map[int]domain.RosstatSlice)
	for i, item := range data {
		if i%1000 == 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
		key := item.SubjectCode
		result[key] = append(result[key], item)
	}
	return result, nil
}

func aggregateByAge(ageGroups [][]*domain.RosstatAge) []*domain.RosstatAge {
	ages := make(map[int]*domain.RosstatAge)
	for _, groups := range ageGroups {
		for _, ageItem := range groups {
			if ageItem == nil {
				continue
			}
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

func buildRegionRecord(subjectCode int, year int, municipalities domain.RosstatSlice) *domain.Rosstat {
	populations := make([]*int, len(municipalities))
	births := make([]*int, len(municipalities))
	deaths := make([]*int, len(municipalities))
	arrivals := make([]*int, len(municipalities))
	departures := make([]*int, len(municipalities))
	males := make([]*int, len(municipalities))
	females := make([]*int, len(municipalities))
	landAreas := make([]*int, len(municipalities))
	avgSalaries := make([]*float64, len(municipalities))
	medicalFacilities := make([]*int, len(municipalities))
	schools := make([]*int, len(municipalities))
	housing := make([]*int, len(municipalities))
	allAgeGroups := make([][]*domain.RosstatAge, len(municipalities))
	for i, m := range municipalities {
		populations[i] = m.Population
		births[i] = m.Birth
		deaths[i] = m.Death
		arrivals[i] = m.Arrival
		departures[i] = m.Departure
		males[i] = m.Male
		females[i] = m.Female
		landAreas[i] = m.LandArea
		avgSalaries[i] = m.AverageSalary
		medicalFacilities[i] = m.MedicalFacilities
		schools[i] = m.Schools
		housing[i] = m.HousingCommissioned
		allAgeGroups[i] = m.ByAge
	}
	return &domain.Rosstat{
		Code:                subjectCode,
		SubjectCode:         subjectCode,
		Year:                year,
		Population:          aggregateInt(populations),
		Birth:               aggregateInt(births),
		Death:               aggregateInt(deaths),
		Arrival:             aggregateInt(arrivals),
		Departure:           aggregateInt(departures),
		Male:                aggregateInt(males),
		Female:              aggregateInt(females),
		LandArea:            aggregateInt(landAreas),
		AverageSalary:       aggregateFloat(avgSalaries),
		MedicalFacilities:   aggregateInt(medicalFacilities),
		Schools:             aggregateInt(schools),
		HousingCommissioned: aggregateInt(housing),
		ByAge:               aggregateByAge(allAgeGroups),
	}
}
