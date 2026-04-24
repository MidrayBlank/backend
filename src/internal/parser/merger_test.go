package parser

import (
	"fmt"
	"testing"
)

func TestFullMerge(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full merge test in short mode")
	}

	collector := NewDataCollector()

	fmt.Println("Starting full data collection...")

	records, err := collector.MergeAll()
	if err != nil {
		t.Fatalf("MergeAll error: %v", err)
	}

	fmt.Printf("\n=== RESULTS ===\n")
	fmt.Printf("Total records: %d\n", len(records))

	// Статистика по заполненности полей
	var (
		hasPopulation, hasBirths, hasDeaths                    int
		hasArrival, hasDeparture                               int
		hasMale, hasFemale                                     int
		hasLand, hasSalary, hasMedical, hasSchools, hasHousing int
	)

	for _, r := range records {
		if r.PopulationAmount > 0 {
			hasPopulation++
		}
		if r.BirthAmount != nil {
			hasBirths++
		}
		if r.DeathAmount != nil {
			hasDeaths++
		}
		if r.ArrivalAmount != nil {
			hasArrival++
		}
		if r.DepartureAmount != nil {
			hasDeparture++
		}
		if r.MaleAmount != nil {
			hasMale++
		}
		if r.FemaleAmount != nil {
			hasFemale++
		}
		if r.LandArea != nil {
			hasLand++
		}
		if r.AvgSalary != nil {
			hasSalary++
		}
		if r.MedicalFacilities != nil {
			hasMedical++
		}
		if r.SchoolsCount != nil {
			hasSchools++
		}
		if r.HousingCommissioned != nil {
			hasHousing++
		}
	}

	fmt.Println("\nField coverage:")
	fmt.Printf("  PopulationAmount:  %d / %d (%.1f%%)\n", hasPopulation, len(records), float64(hasPopulation)/float64(len(records))*100)
	fmt.Printf("  BirthAmount:       %d / %d (%.1f%%)\n", hasBirths, len(records), float64(hasBirths)/float64(len(records))*100)
	fmt.Printf("  DeathAmount:       %d / %d (%.1f%%)\n", hasDeaths, len(records), float64(hasDeaths)/float64(len(records))*100)
	fmt.Printf("  ArrivalAmount:     %d / %d (%.1f%%)\n", hasArrival, len(records), float64(hasArrival)/float64(len(records))*100)
	fmt.Printf("  DepartureAmount:   %d / %d (%.1f%%)\n", hasDeparture, len(records), float64(hasDeparture)/float64(len(records))*100)
	fmt.Printf("  MaleAmount:        %d / %d (%.1f%%)\n", hasMale, len(records), float64(hasMale)/float64(len(records))*100)
	fmt.Printf("  FemaleAmount:      %d / %d (%.1f%%)\n", hasFemale, len(records), float64(hasFemale)/float64(len(records))*100)
	fmt.Printf("  LandArea:          %d / %d (%.1f%%)\n", hasLand, len(records), float64(hasLand)/float64(len(records))*100)
	fmt.Printf("  AvgSalary:         %d / %d (%.1f%%)\n", hasSalary, len(records), float64(hasSalary)/float64(len(records))*100)
	fmt.Printf("  MedicalFacilities: %d / %d (%.1f%%)\n", hasMedical, len(records), float64(hasMedical)/float64(len(records))*100)
	fmt.Printf("  SchoolsCount:      %d / %d (%.1f%%)\n", hasSchools, len(records), float64(hasSchools)/float64(len(records))*100)
	fmt.Printf("  HousingCommissioned: %d / %d (%.1f%%)\n", hasHousing, len(records), float64(hasHousing)/float64(len(records))*100)

	// Примеры записей
	fmt.Println("\nSample records (first 5):")
	for i, r := range records {
		if i >= 5 {
			break
		}
		fmt.Printf("\n  %d: code=%s, year=%d, population=%d\n", i+1, r.Code, r.Year, r.PopulationAmount)
		if r.BirthAmount != nil {
			fmt.Printf("      births=%d", r.BirthAmount)
		}
		if r.DeathAmount != nil {
			fmt.Printf(", deaths=%d", r.DeathAmount)
		}
		if r.ArrivalAmount != nil {
			fmt.Printf(", arrival=%d", r.ArrivalAmount)
		}
		if r.DepartureAmount != nil {
			fmt.Printf(", departure=%d", r.DepartureAmount)
		}
		fmt.Println()
	}
}

// Быстрый тест для проверки структуры
func TestStructureOnly(t *testing.T) {
	record := DemographyRecord{
		Code:             "12345678901",
		Year:             2022,
		PopulationAmount: 10000,
	}

	if record.Code == "" {
		t.Error("Code should not be empty")
	}
	if record.Year == 0 {
		t.Error("Year should not be 0")
	}
	if record.PopulationAmount == 0 {
		t.Error("PopulationAmount should not be 0")
	}

	t.Log("Structure test passed")
}
