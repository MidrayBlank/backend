package impl

import (
	"backend/src/internal/domain"
	"testing"
)

func TestToRosstatList(t *testing.T) {
	service := &RosstatService{}

	t.Run("successful combinating", func(t *testing.T) {
		population1 := 1200000
		birth1 := 12000
		death1 := 15000
		arrival1 := 5000
		departure1 := 3000
		male1 := 560000
		female1 := 640000
		landArea1 := 500
		avgSalary1 := 55000.0
		medical1 := 45
		schools1 := 120
		housing1 := 250

		population2 := 1210000
		birth2 := 11800
		death2 := 15200
		arrival2 := 5200
		departure2 := 3100
		male2 := 565000
		female2 := 645000
		landArea2 := 505
		avgSalary2 := 56000.0
		medical2 := 46
		schools2 := 121
		housing2 := 260

		rosstatByYears := []*domain.Rosstat{
			{
				ID:                  1,
				Code:                45001,
				Year:                2023,
				PopulationAmount:    &population1,
				BirthAmount:         &birth1,
				DeathAmount:         &death1,
				ArrivalAmount:       &arrival1,
				DepartureAmount:     &departure1,
				MaleAmount:          &male1,
				FemaleAmount:        &female1,
				LandArea:            &landArea1,
				AvgSalary:           &avgSalary1,
				MedicalFacilities:   &medical1,
				SchoolsCount:        &schools1,
				HousingCommissioned: &housing1,
			},
			{
				ID:                  2,
				Code:                45001,
				Year:                2024,
				PopulationAmount:    &population2,
				BirthAmount:         &birth2,
				DeathAmount:         &death2,
				ArrivalAmount:       &arrival2,
				DepartureAmount:     &departure2,
				MaleAmount:          &male2,
				FemaleAmount:        &female2,
				LandArea:            &landArea2,
				AvgSalary:           &avgSalary2,
				MedicalFacilities:   &medical2,
				SchoolsCount:        &schools2,
				HousingCommissioned: &housing2,
			},
		}

		rosstatAges := []*domain.RosstatByAge{
			{RosstatID: 1, Age: 0, MaleAmount: 5000, FemaleAmount: 4800},
			{RosstatID: 1, Age: 25, MaleAmount: 7000, FemaleAmount: 6800},
			{RosstatID: 2, Age: 0, MaleAmount: 5100, FemaleAmount: 4900},
		}

		result, err := service.toRosstatList(rosstatByYears, rosstatAges)
		t.Log(result)

		if err != nil {
			t.Errorf("expected no errors")
		}

		if result == nil {
			t.Errorf("expected not nill result")
		}

		if len(result) != 2 {
			t.Errorf("invalid count of objects")
		}

		if *result[0].HousingCommissioned != 250 {
			t.Errorf("got: %d, expected: %f", *result[0].HousingCommissioned, 250.3)
		}

	})
}
