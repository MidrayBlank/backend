package yadisk

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"backend/pkg/parser/rosstat/abstract"
	"backend/pkg/parser/rosstat/downloader"
	"backend/pkg/parser/rosstat/merge"
	"backend/pkg/parser/rosstat/parser"
)

// URL переменные
var (
	populationURL = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112027_112_v20250918.zip"
	birthsURL     = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112003_112_v20250918.zip"
	deathsURL     = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112001_112_v20250918.zip"
	landURL       = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section6/data_Y48006001_112_v20250918.zip"
	healthcareURL = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section33/data_Y48018000_112_v20250918.zip"
	educationURL  = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section34/data_Y48015001_112_v20250918.zip"
	housingURL    = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section37/data_Y48010001_112_v20250918.zip"
	salaryURL     = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section32/data_Y48423007_112_v20250918.zip"
	ageSexURL     = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112014_112_v20250918.zip"
	arrivalURL    = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112021_112_v20250918.zip"
	departureURL  = "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112022_112_v20250918.zip"
)

type YadiskParser struct{}

func NewYadiskParser() *YadiskParser {
	return &YadiskParser{}
}

func (p *YadiskParser) Parse(ctx context.Context) (*abstract.ParseResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Создаём менеджер карты записей
	manager := merge.NewRecordMapManager()

	// 1. Загружаем население
	if err := p.loadPopulation(ctx, manager); err != nil {
		return nil, fmt.Errorf("population load failed: %w", err)
	}

	// 2. Загружаем рождения и смерти
	if err := p.loadBirths(ctx, manager); err != nil {
		return nil, fmt.Errorf("births load failed: %w", err)
	}
	if err := p.loadDeaths(ctx, manager); err != nil {
		return nil, fmt.Errorf("deaths load failed: %w", err)
	}

	// 3. Параллельная загрузка остальных показателей
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return p.loadMultipleParallel(ctx, manager)
	})

	g.Go(func() error {
		return p.loadArrival(ctx, manager)
	})

	g.Go(func() error {
		return p.loadDeparture(ctx, manager)
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("parallel loading failed: %w", err)
	}

	// 4. Проверка контекста
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// 5. Возвращаем результат
	return &abstract.ParseResult{
		Records: manager.GetRecords(),
	}, nil
}

// loadPopulation загружает численность населения
func (p *YadiskParser) loadPopulation(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, populationURL)
	if err != nil {
		return fmt.Errorf("population download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	populationParser := parser.NewPopulationParser()
	records, err := populationParser.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("population parse error: %w", err)
	}

	for _, r := range records {
		manager.SetPopulation(r.Oktmo, r.Year, int(r.Value))
	}
	fmt.Printf("  Population: %d records\n", len(records))
	return nil
}

// loadBirths загружает число родившихся
func (p *YadiskParser) loadBirths(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, birthsURL)
	if err != nil {
		return fmt.Errorf("births download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	birthDeathParser := parser.NewBirthDeathParser()
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
func (p *YadiskParser) loadDeaths(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, deathsURL)
	if err != nil {
		return fmt.Errorf("deaths download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	birthDeathParser := parser.NewBirthDeathParser()
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

// loadMultipleParallel загружает несколько показателей параллельно
func (p *YadiskParser) loadMultipleParallel(ctx context.Context, manager *merge.RecordMapManager) error {
	type task struct {
		name string
		fn   func(context.Context, *merge.RecordMapManager) error
	}

	tasks := []task{
		{"land_area", p.loadLandArea},
		{"healthcare", p.loadHealthcare},
		{"education", p.loadEducation},
		{"housing", p.loadHousing},
		{"salary", p.loadSalary},
		{"age_sex", p.loadAgeSex},
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, t := range tasks {
		task := t
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			fmt.Printf("  Loading %s...\n", task.name)
			return task.fn(ctx, manager)
		})
	}

	return g.Wait()
}

// loadLandArea загружает площадь территории
func (p *YadiskParser) loadLandArea(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, landURL)
	if err != nil {
		return fmt.Errorf("land download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	landParser := parser.NewLandParser()
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
func (p *YadiskParser) loadHealthcare(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, healthcareURL)
	if err != nil {
		return fmt.Errorf("healthcare download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	healthcareParser := parser.NewHealthcareParser()
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
func (p *YadiskParser) loadEducation(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, educationURL)
	if err != nil {
		return fmt.Errorf("education download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	educationParser := parser.NewEducationParser()
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
func (p *YadiskParser) loadHousing(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, housingURL)
	if err != nil {
		return fmt.Errorf("housing download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	housingParser := parser.NewHousingParser()
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
func (p *YadiskParser) loadSalary(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, salaryURL)
	if err != nil {
		return fmt.Errorf("salary download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	salaryParser := parser.NewSalaryParser()
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
func (p *YadiskParser) loadAgeSex(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := downloader.DownloadAndUnzip(ctx, ageSexURL)
	if err != nil {
		return fmt.Errorf("agesex download error: %w", err)
	}
	defer downloader.CleanupTemp(file)

	if err := ctx.Err(); err != nil {
		return err
	}

	ageSexParser := parser.NewAgeSexParser()
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
func (p *YadiskParser) loadArrival(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fmt.Println("Loading arrival data...")
	files, err := downloader.DownloadAndUnzipAll(ctx, arrivalURL)
	if err != nil {
		return fmt.Errorf("arrival download error: %w", err)
	}
	defer downloader.CleanupTempAll(filepath.Dir(files[0]))

	migrationParser := parser.NewMigrationParser()
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
func (p *YadiskParser) loadDeparture(ctx context.Context, manager *merge.RecordMapManager) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fmt.Println("Loading departure data...")
	files, err := downloader.DownloadAndUnzipAll(ctx, departureURL)
	if err != nil {
		return fmt.Errorf("departure download error: %w", err)
	}
	defer downloader.CleanupTempAll(filepath.Dir(files[0]))

	migrationParser := parser.NewMigrationParser()
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
