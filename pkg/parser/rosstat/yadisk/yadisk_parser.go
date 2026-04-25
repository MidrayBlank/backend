package yadisk

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"backend/pkg/parser/rosstat/yadisk/config"
	"backend/pkg/parser/rosstat/yadisk/downloader"
	"backend/pkg/parser/rosstat/yadisk/extractor"
	"backend/pkg/parser/rosstat/yadisk/merge"
	"backend/pkg/parser/rosstat/yadisk/model"

	"golang.org/x/sync/errgroup"
)

type YadiskRosstatParser struct{}

func NewYadiskRosstatParser() *YadiskRosstatParser {
	return &YadiskRosstatParser{}
}

func (p *YadiskRosstatParser) Parse(ctx context.Context) (model.RosstatParsedSlice, error) {
	cfg := config.NewConfig()
	manager := merge.NewRecordMapManager()

	if err := p.loadPopulation(ctx, manager, cfg.PopulationURL); err != nil {
		return nil, fmt.Errorf("population load failed: %w", err)
	}

	if err := p.loadBirths(ctx, manager, cfg.BirthURL); err != nil {
		return nil, fmt.Errorf("births load failed: %w", err)
	}
	if err := p.loadDeaths(ctx, manager, cfg.DeathURL); err != nil {
		return nil, fmt.Errorf("deaths load failed: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return p.loadLandArea(ctx, manager, cfg.LandAreaURL)
	})

	g.Go(func() error {
		return p.loadHealthcare(ctx, manager, cfg.MedicialFacilitiesURL)
	})

	g.Go(func() error {
		return p.loadEducation(ctx, manager, cfg.SchoolsURL)
	})

	g.Go(func() error {
		return p.loadHousing(ctx, manager, cfg.HousingCommissionedURL)
	})

	g.Go(func() error {
		return p.loadSalary(ctx, manager, cfg.AverageSalaryURL)
	})

	g.Go(func() error {
		return p.loadAgeSex(ctx, manager, cfg.MaleFemaleAgeURL)
	})

	g.Go(func() error {
		return p.loadArrival(ctx, manager, cfg.ArrivalURL)
	})

	g.Go(func() error {
		return p.loadDeparture(ctx, manager, cfg.DepartureURL)
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("parallel loading failed: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return manager.GetRecords(), nil
}

func (p *YadiskRosstatParser) loadPopulation(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("population download error: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	populationParser := extractor.NewPopulationParser()
	records, err := populationParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("population parse error: %w", err)
	}
	downloader.CleanupTemp(file)

	for _, r := range records {
		manager.SetPopulation(r.Oktmo, r.Year, int(r.Value))
	}
	fmt.Printf("  Population: %d records\n", len(records))
	return nil
}

// loadBirths загружает число родившихся
func (p *YadiskRosstatParser) loadBirths(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("births download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	birthDeathParser := extractor.NewBirthDeathParser()
	records, err := birthDeathParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("births parse error: %w", err)
	}

	for _, r := range records {
		manager.SetBirths(r.Oktmo, r.Year, int(r.Value))
	}
	fmt.Printf("  Births: %d records\n", len(records))
	return nil
}

// loadDeaths загружает число умерших
func (p *YadiskRosstatParser) loadDeaths(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("deaths download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	birthDeathParser := extractor.NewBirthDeathParser()
	records, err := birthDeathParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("deaths parse error: %w", err)
	}

	for _, r := range records {
		manager.SetDeaths(r.Oktmo, r.Year, int(r.Value))
	}
	fmt.Printf("  Deaths: %d records\n", len(records))
	return nil
}

// loadLandArea загружает площадь территории
func (p *YadiskRosstatParser) loadLandArea(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("land download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	landParser := extractor.NewLandParser()
	records, err := landParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("land parse error: %w", err)
	}

	aggregated := downloader.AggregateByDistrict(ctx, &records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetLandArea(oktmo, year, value)
		}
	}
	fmt.Printf("  Land area: %d records aggregated\n", len(records))
	return nil
}

// loadHealthcare загружает количество медучреждений
func (p *YadiskRosstatParser) loadHealthcare(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("healthcare download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	healthcareParser := extractor.NewHealthcareParser()
	records, err := healthcareParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("healthcare parse error: %w", err)
	}

	aggregated := downloader.AggregateByDistrict(ctx, &records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetMedicalFacilities(oktmo, year, int(value))
		}
	}
	fmt.Printf("  Healthcare: %d records aggregated\n", len(records))
	return nil
}

// loadEducation загружает количество школ
func (p *YadiskRosstatParser) loadEducation(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("education download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	educationParser := extractor.NewEducationParser()
	records, err := educationParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("education parse error: %w", err)
	}

	aggregated := downloader.AggregateByDistrict(ctx, &records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetSchoolsCount(oktmo, year, int(value))
		}
	}
	fmt.Printf("  Education: %d records aggregated\n", len(records))
	return nil
}

// loadHousing загружает введённое жильё
func (p *YadiskRosstatParser) loadHousing(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("housing download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	housingParser := extractor.NewHousingParser()
	records, err := housingParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("housing parse error: %w", err)
	}

	aggregated := downloader.AggregateSimple(ctx, &records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetHousingCommissioned(oktmo, year, value)
		}
	}
	fmt.Printf("  Housing: %d records aggregated\n", len(records))
	return nil
}

// loadSalary загружает среднюю зарплату
func (p *YadiskRosstatParser) loadSalary(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("salary download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	salaryParser := extractor.NewSalaryParser()
	records, err := salaryParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("salary parse error: %w", err)
	}

	aggregated := downloader.AggregateSimple(ctx, &records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetAvgSalary(oktmo, year, value)
		}
	}
	fmt.Printf("  Salary: %d records aggregated\n", len(records))
	return nil
}

// loadAgeSex загружает половозрастную структуру
func (p *YadiskRosstatParser) loadAgeSex(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, url)
	if err != nil {
		return fmt.Errorf("agesex download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	ageSexParser := extractor.NewAgeSexParser()
	records, err := ageSexParser.ParseAgeSex(ctx, file)
	if err != nil {
		return fmt.Errorf("agesex parse error: %w", err)
	}

	// Агрегируем для основной таблицы (суммарные male/female)
	aggregated := make(map[string]map[int]struct{ male, female int })
	for _, r := range records {
		normalizedOktmo := downloader.NormalizeOktmo(r.Oktmo)
		if aggregated[normalizedOktmo] == nil {
			aggregated[normalizedOktmo] = make(map[int]struct{ male, female int })
		}
		entry := aggregated[normalizedOktmo][r.Year]
		if r.Sex == "Мужчины" {
			entry.male += r.Value
		} else {
			entry.female += r.Value
		}
		aggregated[normalizedOktmo][r.Year] = entry
	}

	for oktmo, years := range aggregated {
		for year, values := range years {
			manager.SetAgeSex(oktmo, year, values.male, values.female)
		}
	}
	fmt.Printf("  AgeSex: %d raw records\n", len(records))
	return nil
}

// loadArrival загружает прибывших
func (p *YadiskRosstatParser) loadArrival(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fmt.Println("Loading arrival data...")
	files, err := downloader.DownloadAndUnzipAll(ctx, url)
	if err != nil {
		return fmt.Errorf("arrival download error: %w", err)
	}
	defer downloader.CleanupTempAll(filepath.Dir(files[0]))

	migrationParser := extractor.NewMigrationParser()
	var allRecords []downloader.MigrationRecord

	for _, file := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if strings.Contains(file, "__MACOSX") || strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}
		records, err := migrationParser.ParseMigration(ctx, file)
		if err != nil {
			log.Printf("Warning: arrival parse error for %s: %v", filepath.Base(file), err)
			continue
		}
		allRecords = append(allRecords, records...)
	}

	aggregated := downloader.AggregateMigrationWithNormalize(ctx, &allRecords)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetArrival(oktmo, year, value)
		}
	}
	fmt.Printf("  Arrival: %d records\n", len(allRecords))
	return nil
}

// loadDeparture загружает убывших
func (p *YadiskRosstatParser) loadDeparture(ctx context.Context, manager *merge.RecordMapManager, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fmt.Println("Loading departure data...")
	files, err := downloader.DownloadAndUnzipAll(ctx, url)
	if err != nil {
		return fmt.Errorf("departure download error: %w", err)
	}
	defer downloader.CleanupTempAll(filepath.Dir(files[0]))

	migrationParser := extractor.NewMigrationParser()
	var allRecords []downloader.MigrationRecord

	for _, file := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if strings.Contains(file, "__MACOSX") || strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}
		records, err := migrationParser.ParseMigration(ctx, file)
		if err != nil {
			log.Printf("Warning: departure parse error for %s: %v", filepath.Base(file), err)
			continue
		}
		allRecords = append(allRecords, records...)
	}

	aggregated := downloader.AggregateMigrationWithNormalize(ctx, &allRecords)
	for oktmo, years := range aggregated {
		for year, value := range years {
			manager.SetDeparture(oktmo, year, value)
		}
	}
	fmt.Printf("  Departure: %d records\n", len(allRecords))
	return nil
}
