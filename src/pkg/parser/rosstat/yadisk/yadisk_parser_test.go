package yadisk_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"backend/src/pkg/parser/rosstat/yadisk"
	"backend/src/pkg/parser/rosstat/yadisk/config"
	"backend/src/pkg/parser/rosstat/yadisk/storage"
	"backend/src/pkg/parser/rosstat/yadisk/subparser"
)

func TestYadiskRosstatParser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full parser test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	parser := yadisk.NewYadiskRosstatParser()

	fmt.Println("Starting YadiskRosstatParser test...")

	result, err := parser.Parse(ctx)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty records, got 0")
	}

	fmt.Printf("\n=== TEST RESULTS ===\n")
	fmt.Printf("Total records: %d\n", len(result))

	var (
		hasPopulation, hasBirths, hasDeaths                    int
		hasArrival, hasDeparture                               int
		hasMale, hasFemale                                     int
		hasLand, hasSalary, hasMedical, hasSchools, hasHousing int
	)

	for _, r := range result {
		if r.Population != nil {
			hasPopulation++
		}

		if r.Birth != nil && *r.Birth > 0 {
			hasBirths++
		}
		if r.Death != nil && *r.Death > 0 {
			hasDeaths++
		}
		if r.Arrival != nil && *r.Arrival > 0 {
			hasArrival++
		}
		if r.Departure != nil && *r.Departure > 0 {
			hasDeparture++
		}
		if r.Male != nil && *r.Male > 0 {
			hasMale++
		}
		if r.Female != nil && *r.Female > 0 {
			hasFemale++
		}
		if r.LandArea != nil && *r.LandArea > 0 {
			hasLand++
		}
		if r.AverageSalary != nil && *r.AverageSalary > 0 {
			hasSalary++
		}
		if r.MedicialFacilities != nil && *r.MedicialFacilities > 0 {
			hasMedical++
		}
		if r.Schools != nil && *r.Schools > 0 {
			hasSchools++
		}
		if r.HousingCommissioned != nil && *r.HousingCommissioned > 0 {
			hasHousing++
		}
	}

	fmt.Println("\nField coverage:")
	fmt.Printf("  Population:   %d / %d (%.1f%%)\n", hasPopulation, len(result), float64(hasPopulation)/float64(len(result))*100)
	fmt.Printf("  Birth:        %d / %d (%.1f%%)\n", hasBirths, len(result), float64(hasBirths)/float64(len(result))*100)
	fmt.Printf("  Death:        %d / %d (%.1f%%)\n", hasDeaths, len(result), float64(hasDeaths)/float64(len(result))*100)
	fmt.Printf("  Arrival:      %d / %d (%.1f%%)\n", hasArrival, len(result), float64(hasArrival)/float64(len(result))*100)
	fmt.Printf("  Departure:    %d / %d (%.1f%%)\n", hasDeparture, len(result), float64(hasDeparture)/float64(len(result))*100)
	fmt.Printf("  Male:         %d / %d (%.1f%%)\n", hasMale, len(result), float64(hasMale)/float64(len(result))*100)
	fmt.Printf("  Female:       %d / %d (%.1f%%)\n", hasFemale, len(result), float64(hasFemale)/float64(len(result))*100)
	fmt.Printf("  LandArea:           %d / %d (%.1f%%)\n", hasLand, len(result), float64(hasLand)/float64(len(result))*100)
	fmt.Printf("  AvgSalary:          %d / %d (%.1f%%)\n", hasSalary, len(result), float64(hasSalary)/float64(len(result))*100)
	fmt.Printf("  MedicalFacilities:  %d / %d (%.1f%%)\n", hasMedical, len(result), float64(hasMedical)/float64(len(result))*100)
	fmt.Printf("  SchoolsCount:       %d / %d (%.1f%%)\n", hasSchools, len(result), float64(hasSchools)/float64(len(result))*100)
	fmt.Printf("  HousingCommissioned:%d / %d (%.1f%%)\n", hasHousing, len(result), float64(hasHousing)/float64(len(result))*100)

	fmt.Println("\nSample records (first 5):")
	for i, r := range result {
		if i >= 5 {
			break
		}
		fmt.Printf("\n  %d: code=%d, year=%d, population=%d\n", i+1, r.Code, r.Year, r.Population)
		if r.Birth != nil && *r.Birth != 0 {
			fmt.Printf("      births=%d", *r.Birth)
		}
		if r.Death != nil && *r.Death != 0 {
			fmt.Printf(", deaths=%d", *r.Death)
		}
		if r.Arrival != nil && *r.Arrival != 0 {
			fmt.Printf(", arrival=%d", *r.Arrival)
		}
		if r.Departure != nil && *r.Departure != 0 {
			fmt.Printf(", departure=%d", *r.Departure)
		}
		fmt.Println()
	}

	if hasPopulation == 0 {
		t.Error("No population data loaded")
	}
	if hasArrival == 0 {
		t.Log("Warning: No arrival data loaded")
	}
	if hasDeparture == 0 {
		t.Log("Warning: No departure data loaded")
	}
}

func TestPopulation(t *testing.T) {
	config := config.NewConfig()
	parser := subparser.PopulationSubparser{}
	storage := storage.NewStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := parser.Parse(ctx, storage, config.PopulationURL)

	if err != nil {
		t.Fatalf("Population subparser error: %v", err)
	}

	result := storage.Result()

	fmt.Printf("✅ Success! Parsed %d records by population subpraser\n", len(result))

	fmt.Println("\nSample records:")
	for i, r := range result {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d: Code=%d, Year=%d, Population=%d\n", i+1, r.Code, r.Year, r.Population)
	}

	if len(result) == 0 {
		t.Error("No records parsed")
	}
}
