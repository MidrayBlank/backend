package yadisk

import (
	"backend/pkg/parser/rosstat/downloader"
	"backend/pkg/parser/rosstat/parser"
	"context"
	"fmt"
	"testing"
	"time"
)

// TestYadiskParser проверяет полную работу парсера
func TestYadiskParser(t *testing.T) {
	// Пропускаем тест в коротком режиме (чтобы не тратить время при быстрых тестах)
	if testing.Short() {
		t.Skip("Skipping full parser test in short mode")
	}

	// Создаём контекст с таймаутом 60 минут
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	// Создаём парсер
	parser := NewYadiskParser()

	fmt.Println("Starting YadiskParser test...")
	fmt.Println("This will take about 30-40 minutes to download all data...")

	// Запускаем парсинг
	result, err := parser.Parse(ctx)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Проверяем, что результат не пустой
	if len(result.Records) == 0 {
		t.Error("Expected non-empty records, got 0")
	}

	// Выводим статистику
	fmt.Printf("\n=== TEST RESULTS ===\n")
	fmt.Printf("Total records: %d\n", len(result.Records))

	// Считаем заполненность полей
	var (
		hasPopulation, hasBirths, hasDeaths                    int
		hasArrival, hasDeparture                               int
		hasMale, hasFemale                                     int
		hasLand, hasSalary, hasMedical, hasSchools, hasHousing int
	)

	for _, r := range result.Records {
		if r.PopulationAmount > 0 {
			hasPopulation++
		}
		if r.BirthAmount != nil && *r.BirthAmount > 0 {
			hasBirths++
		}
		if r.DeathAmount != nil && *r.DeathAmount > 0 {
			hasDeaths++
		}
		if r.ArrivalAmount != nil && *r.ArrivalAmount > 0 {
			hasArrival++
		}
		if r.DepartureAmount != nil && *r.DepartureAmount > 0 {
			hasDeparture++
		}
		if r.MaleAmount != nil && *r.MaleAmount > 0 {
			hasMale++
		}
		if r.FemaleAmount != nil && *r.FemaleAmount > 0 {
			hasFemale++
		}
		if r.LandArea != nil && *r.LandArea > 0 {
			hasLand++
		}
		if r.AvgSalary != nil && *r.AvgSalary > 0 {
			hasSalary++
		}
		if r.MedicalFacilities != nil && *r.MedicalFacilities > 0 {
			hasMedical++
		}
		if r.SchoolsCount != nil && *r.SchoolsCount > 0 {
			hasSchools++
		}
		if r.HousingCommissioned != nil && *r.HousingCommissioned > 0 {
			hasHousing++
		}
	}

	fmt.Println("\nField coverage:")
	fmt.Printf("  PopulationAmount:   %d / %d (%.1f%%)\n", hasPopulation, len(result.Records), float64(hasPopulation)/float64(len(result.Records))*100)
	fmt.Printf("  BirthAmount:        %d / %d (%.1f%%)\n", hasBirths, len(result.Records), float64(hasBirths)/float64(len(result.Records))*100)
	fmt.Printf("  DeathAmount:        %d / %d (%.1f%%)\n", hasDeaths, len(result.Records), float64(hasDeaths)/float64(len(result.Records))*100)
	fmt.Printf("  ArrivalAmount:      %d / %d (%.1f%%)\n", hasArrival, len(result.Records), float64(hasArrival)/float64(len(result.Records))*100)
	fmt.Printf("  DepartureAmount:    %d / %d (%.1f%%)\n", hasDeparture, len(result.Records), float64(hasDeparture)/float64(len(result.Records))*100)
	fmt.Printf("  MaleAmount:         %d / %d (%.1f%%)\n", hasMale, len(result.Records), float64(hasMale)/float64(len(result.Records))*100)
	fmt.Printf("  FemaleAmount:       %d / %d (%.1f%%)\n", hasFemale, len(result.Records), float64(hasFemale)/float64(len(result.Records))*100)
	fmt.Printf("  LandArea:           %d / %d (%.1f%%)\n", hasLand, len(result.Records), float64(hasLand)/float64(len(result.Records))*100)
	fmt.Printf("  AvgSalary:          %d / %d (%.1f%%)\n", hasSalary, len(result.Records), float64(hasSalary)/float64(len(result.Records))*100)
	fmt.Printf("  MedicalFacilities:  %d / %d (%.1f%%)\n", hasMedical, len(result.Records), float64(hasMedical)/float64(len(result.Records))*100)
	fmt.Printf("  SchoolsCount:       %d / %d (%.1f%%)\n", hasSchools, len(result.Records), float64(hasSchools)/float64(len(result.Records))*100)
	fmt.Printf("  HousingCommissioned:%d / %d (%.1f%%)\n", hasHousing, len(result.Records), float64(hasHousing)/float64(len(result.Records))*100)

	// Показываем примеры
	fmt.Println("\nSample records (first 5):")
	for i, r := range result.Records {
		if i >= 5 {
			break
		}
		fmt.Printf("\n  %d: code=%s, year=%d, population=%d\n", i+1, r.Code, r.Year, r.PopulationAmount)
		if r.BirthAmount != nil && *r.BirthAmount != 0 {
			fmt.Printf("      births=%d", *r.BirthAmount)
		}
		if r.DeathAmount != nil && *r.DeathAmount != 0 {
			fmt.Printf(", deaths=%d", *r.DeathAmount)
		}
		if r.ArrivalAmount != nil && *r.ArrivalAmount != 0 {
			fmt.Printf(", arrival=%d", *r.ArrivalAmount)
		}
		if r.DepartureAmount != nil && *r.DepartureAmount != 0 {
			fmt.Printf(", departure=%d", *r.DepartureAmount)
		}
		fmt.Println()
	}

	// Проверяем, что есть данные
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

func TestQuickPopulation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Println("Testing population download and parse...")

	// Скачиваем файл населения
	file, err := downloader.DownloadAndUnzip(ctx, populationURL)
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}
	defer downloader.CleanupTemp(file)

	// Парсим
	populationParser := parser.NewPopulationParser()
	records, err := populationParser.Parse(ctx, file)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	fmt.Printf("✅ Success! Parsed %d population records\n", len(records))

	// Показываем первые 5 записей
	fmt.Println("\nSample records:")
	for i, r := range records {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d: OKTMO=%s, Year=%d, Population=%.0f\n", i+1, r.Oktmo, r.Year, r.Value)
	}

	if len(records) == 0 {
		t.Error("No records parsed")
	}
}
