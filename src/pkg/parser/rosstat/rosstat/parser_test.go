package rosstat_test

import (
	"context"
	"testing"
	"time"

	"backend/src/pkg/parser/rosstat/rosstat/config"
	"backend/src/pkg/parser/rosstat/rosstat/storage"
	"backend/src/pkg/parser/rosstat/rosstat/subparser"
)

func TestRosstatParserPopulation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full parser test in short mode")
	}

	storage := storage.NewStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	code := subparser.NewCodeSubparser()
	err := code.Parse(ctx, storage)
	if err != nil {
		t.Errorf("Error occurred in code parser: %s", err.Error())
	}

	population := subparser.NewPopulationSubparser()

	err = population.Parse(ctx, storage)
	if err != nil {
		t.Errorf("Error occured: %s", err.Error())
	}

	config := config.NewConfig()
	subjectCodes := config.SubjectCodes

	result := storage.Result()

	subjectCodesBool := make([]bool, 100)
	for _, item := range result {
		if item.ParentCode < 100 {
			subjectCodesBool[item.ParentCode] = true
		}

		// t.Log(item.ParentCode, item.Code, item.Year, *item.Population)
	}

	for _, subjectCode := range subjectCodes {
		if !subjectCodesBool[subjectCode] {
			t.Log("SubjectCode ", subjectCode, " not parsed")
		}
	}
}

func TestRosstatParserBirth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full parser test in short mode")
	}

	storage := storage.NewStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()

	code := subparser.NewCodeSubparser()
	err := code.Parse(ctx, storage)
	if err != nil {
		t.Errorf("Error occurred in code parser: %s", err.Error())
	}

	birth := subparser.NewBirthSubparser()
	err = birth.Parse(ctx, storage)
	if err != nil {
		t.Errorf("Error occurred in birth parser: %s", err.Error())
	}

	config := config.NewConfig()
	subjectCodes := config.SubjectCodes

	result := storage.Result()

	subjectCodesBool := make([]bool, 100)
	for _, item := range result {
		subjectCodesBool[getSubjectCode(item.Code)] = true
	}

	for _, subjectCode := range subjectCodes {
		if !subjectCodesBool[subjectCode] {
			t.Log("SubjectCode ", subjectCode, " not parsed")
		}
	}
}

func getSubjectCode(code int) int {
	subjectCode := code / 1000000
	return subjectCode
}

func TestRosstatParserDeath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full parser test in short mode")
	}

	storage := storage.NewStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	code := subparser.NewCodeSubparser()
	err := code.Parse(ctx, storage)
	if err != nil {
		t.Errorf("Error occurred in code parser: %s", err.Error())
	}

	death := subparser.NewDeathSubparser()
	err = death.Parse(ctx, storage)
	if err != nil {
		t.Errorf("Error occurred in death parser: %s", err.Error())
	}

	result := storage.Result()
	deathCount := 0
	for _, item := range result {
		if item.Death != nil {
			deathCount++
		}
	}
	t.Logf("Records with death data: %d", deathCount)
}
