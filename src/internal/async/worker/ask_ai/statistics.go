package ask_ai

import (
	"backend/src/internal/async/worker/repository"
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"context"
	"fmt"
	"strings"
)

type AgeGroup struct {
	Young  int
	Adult  int
	Senior int
}

func prepareStatistics(
	ctx context.Context,
	conn abstract.IDBConnection,
	repositories repository.WorkerRepositories,
	code int,
) ([]*domain.Rosstat, []*domain.RosstatByAge, error) {
	rosstat, err := repositories.RosstatRepository.GetRosstatByCodeForLastFiveYears(conn, code)
	if err != nil {
		return nil, nil, err
	}

	rosstatIDs := make([]int, 0, len(rosstat))
	for _, r := range rosstat {
		rosstatIDs = append(rosstatIDs, r.ID)
	}

	rosstatAge, err := repositories.RosstatAgeRepository.GetRosstatAgeByRosstatIDs(conn, rosstatIDs)
	if err != nil {
		return nil, nil, err
	}

	return rosstat, rosstatAge, nil
}

func buildAIPrompt(regionName string, stats string) string {
	prompt := fmt.Sprintf(AiRequestTemplate, regionName, stats)

	return prompt
}

func formatStatsForAI(stats []*domain.Rosstat, ageStats []*domain.RosstatByAge) string {
	var sb strings.Builder
	sb.WriteString(tableHeader)
	sb.WriteString(tableSeparator)

	ageGrouped := groupPopulationByAge(ageStats)

	for _, s := range stats {
		population := formatIntValue(s.PopulationAmount)
		birth := formatIntValue(s.BirthAmount)
		death := formatIntValue(s.DeathAmount)
		arrival := formatIntValue(s.ArrivalAmount)
		departure := formatIntValue(s.DepartureAmount)
		avgSalary := formatFloatValue(s.AvgSalary)
		male := formatIntValue(s.MaleAmount)
		female := formatIntValue(s.FemaleAmount)
		schools := formatIntValue(s.SchoolsCount)
		medical := formatIntValue(s.MedicalFacilities)
		housing := formatIntValue(s.HousingCommissioned)

		young, adult, senior := getAgeGroups(ageGrouped, s.ID)
		youngStr := formatIntValue(young)
		adultStr := formatIntValue(adult)
		seniorStr := formatIntValue(senior)

		sb.WriteString(fmt.Sprintf(tableRowFormat,
			s.Year,
			population, birth, death,
			arrival, departure,
			avgSalary,
			male, female,
			schools, medical, housing,
			youngStr, adultStr, seniorStr))
	}

	return sb.String()
}

func groupPopulationByAge(ageStats []*domain.RosstatByAge) map[int]*AgeGroup {
	groupedPopulation := make(map[int]*AgeGroup)

	for _, a := range ageStats {
		key := a.RosstatID

		if a.Age < 18 {
			groupedPopulation[key].Young += a.MaleAmount + a.FemaleAmount
		} else if a.Age <= 60 {
			groupedPopulation[key].Adult += a.MaleAmount + a.FemaleAmount
		} else {
			groupedPopulation[key].Senior += a.MaleAmount + a.FemaleAmount
		}
	}

	return groupedPopulation
}

func getAgeGroups(groupedPopulation map[int]*AgeGroup, rosstatID int) (young, adult, senior *int) {
	if group, ok := groupedPopulation[rosstatID]; ok {
		young = &group.Young
		adult = &group.Adult
		senior = &group.Senior
		return
	}
	return nil, nil, nil
}

func formatIntValue(value *int) string {
	if value == nil {
		return "—"
	}
	return fmt.Sprintf("%d", *value)
}

func formatFloatValue(value *float64) string {
	if value == nil {
		return "—"
	}
	return fmt.Sprintf("%f", *value)
}
