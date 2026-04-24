package parser

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DataCollector struct {
	manager *ParserManager
}

func NewDataCollector() *DataCollector {
	return &DataCollector{
		manager: NewParserManager(),
	}
}

// MergeAll собирает все данные и возвращает список записей для БД
func (c *DataCollector) MergeAll() ([]DemographyRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	demRecords, err := c.loadDemography(ctx)
	if err != nil {
		return nil, fmt.Errorf("demography error: %w", err)
	}

	recordMap := c.buildRecordMap(demRecords)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.loadMultipleParallel(ctx, recordMap)
	}()

	wg.Add(1)
	go c.loadArrivalParallel(ctx, recordMap, &wg)

	wg.Add(1)
	go c.loadDepartureParallel(ctx, recordMap, &wg)

	wg.Wait()

	return c.convertToDBRecords(recordMap), nil
}

func (c *DataCollector) loadDemography(ctx context.Context) ([]DemographyRawRecord, error) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_mun_demography_123_v20240612/data_mun_demography_123_v20240612_csv.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		return nil, err
	}
	defer CleanupTemp(file)

	records, err := c.manager.Demography.ParseDemography(file)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c *DataCollector) buildRecordMap(records []DemographyRawRecord) map[string]*DemographyRecord {
	recordMap := make(map[string]*DemographyRecord)

	for _, r := range records {
		key := fmt.Sprintf("%s_%d", r.Oktmo, r.Year)

		recordMap[key] = &DemographyRecord{
			Code:             r.Oktmo,
			Year:             r.Year,
			PopulationAmount: r.Population,
			BirthAmount:      IntPtr(r.Births),
			DeathAmount:      IntPtr(r.Deaths),
		}
	}

	return recordMap
}

func (c *DataCollector) loadLandArea(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section6/data_Y48006001_112_v20250918.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		log.Printf("Warning: land download error: %v", err)
		return
	}
	defer CleanupTemp(file)

	records, err := c.manager.Rosstat.ParseRosstat(file)
	if err != nil {
		log.Printf("Warning: land parse error: %v", err)
		return
	}

	aggregated := AggregateByDistrict(&records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.LandArea = Float64Ptr(value)
			}
		}
	}
}

func (c *DataCollector) loadHealthcare(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section33/data_Y48018000_112_v20250918.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		log.Printf("Warning: healthcare download error: %v", err)
		return
	}
	defer CleanupTemp(file)

	records, err := c.manager.Rosstat.ParseRosstat(file)
	if err != nil {
		log.Printf("Warning: healthcare parse error: %v", err)
		return
	}

	aggregated := AggregateByDistrict(&records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.MedicalFacilities = IntPtr(int(value))
			}
		}
	}
}

func (c *DataCollector) loadEducation(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section34/data_Y48015001_112_v20250918.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		log.Printf("Warning: education download error: %v", err)
		return
	}
	defer CleanupTemp(file)

	records, err := c.manager.Rosstat.ParseRosstat(file)
	if err != nil {
		log.Printf("Warning: education parse error: %v", err)
		return
	}

	aggregated := AggregateByDistrict(&records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.SchoolsCount = IntPtr(int(value))
			}
		}
	}
}

func (c *DataCollector) loadHousing(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section37/data_Y48010001_112_v20250918.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		log.Printf("Warning: housing download error: %v", err)
		return
	}
	defer CleanupTemp(file)

	records, err := c.manager.Rosstat.ParseRosstat(file)
	if err != nil {
		log.Printf("Warning: housing parse error: %v", err)
		return
	}

	aggregated := AggregateSimple(&records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.HousingCommissioned = Float64Ptr(value)
			}
		}
	}
}

func (c *DataCollector) loadSalary(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section32/data_Y48423007_112_v20250918.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		log.Printf("Warning: salary download error: %v", err)
		return
	}
	defer CleanupTemp(file)

	records, err := c.manager.Rosstat.ParseRosstat(file)
	if err != nil {
		log.Printf("Warning: salary parse error: %v", err)
		return
	}

	aggregated := AggregateSimple(&records)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.AvgSalary = Float64Ptr(value)
			}
		}
	}
}

func (c *DataCollector) loadAgeSex(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112014_112_v20250918.zip"

	file, err := DownloadAndUnzip(ctx, url)
	if err != nil {
		log.Printf("Warning: agesex download error: %v", err)
		return
	}
	defer CleanupTemp(file)

	records, err := c.manager.AgeSex.ParseAgeSex(file)
	if err != nil {
		log.Printf("Warning: agesex parse error: %v", err)
		return
	}

	aggregated := AggregateAgeSexWithNormalize(&records)

	for oktmo, years := range aggregated {
		for year, values := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.MaleAmount = IntPtr(values.Male)
				d.FemaleAmount = IntPtr(values.Female)
			}
		}
	}
}

func (c *DataCollector) loadArrival(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112021_112_v20250918.zip"

	fmt.Println("Loading arrival data...")
	files, err := DownloadAndUnzipAll(ctx, url)
	if err != nil {
		log.Printf("Warning: arrival download error: %v", err)
		return
	}
	defer CleanupTempAll(filepath.Dir(files[0]))

	var allRecords []MigrationRecord
	for _, file := range files {
		if strings.Contains(file, "__MACOSX") || strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}
		records, err := c.manager.Migration.ParseMigration(file)
		if err != nil {
			log.Printf("Warning: arrival parse error for %s: %v", filepath.Base(file), err)
			continue
		}
		allRecords = append(allRecords, records...)
	}

	aggregated := AggregateMigrationWithNormalize(&allRecords)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.ArrivalAmount = IntPtr(value)
			}
		}
	}
	fmt.Printf("  Arrival: %d records\n", len(allRecords))
}

func (c *DataCollector) loadDeparture(ctx context.Context, recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112022_112_v20250918.zip"

	fmt.Println("Loading departure data...")
	files, err := DownloadAndUnzipAll(ctx, url)
	if err != nil {
		log.Printf("Warning: departure download error: %v", err)
		return
	}
	defer CleanupTempAll(filepath.Dir(files[0]))

	var allRecords []MigrationRecord
	for _, file := range files {
		if strings.Contains(file, "__MACOSX") || strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}
		records, err := c.manager.Migration.ParseMigration(file)
		if err != nil {
			log.Printf("Warning: departure parse error for %s: %v", filepath.Base(file), err)
			continue
		}
		allRecords = append(allRecords, records...)
	}

	aggregated := AggregateMigrationWithNormalize(&allRecords)
	for oktmo, years := range aggregated {
		for year, value := range years {
			key := fmt.Sprintf("%s_%d", oktmo, year)
			if d, ok := recordMap[key]; ok {
				d.DepartureAmount = IntPtr(value)
			}
		}
	}
	fmt.Printf("  Departure: %d records\n", len(allRecords))
}

func (c *DataCollector) convertToDBRecords(recordMap map[string]*DemographyRecord) []DemographyRecord {
	result := make([]DemographyRecord, 0, len(recordMap))
	for _, record := range recordMap {
		result = append(result, *record)
	}
	return result
}

// loadMultipleParallel загружает несколько показателей параллельно
func (c *DataCollector) loadMultipleParallel(ctx context.Context, recordMap map[string]*DemographyRecord) {
	type task struct {
		name string
		fn   func(context.Context, map[string]*DemographyRecord)
	}

	tasks := []task{
		{"land_area", c.loadLandArea},
		{"healthcare", c.loadHealthcare},
		{"education", c.loadEducation},
		{"housing", c.loadHousing},
		{"salary", c.loadSalary},
		{"age_sex", c.loadAgeSex},
	}

	tasksChan := make(chan task, len(tasks))
	for _, t := range tasks {
		tasksChan <- t
	}
	close(tasksChan)

	var wg sync.WaitGroup
	maxWorkers := 3

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for t := range tasksChan {
				fmt.Printf("  Worker %d: loading %s...\n", workerID, t.name)
				t.fn(ctx, recordMap)
			}
		}(i)
	}

	wg.Wait()
}

// loadArrivalParallel - загружает прибывших (отдельная горутина)
func (c *DataCollector) loadArrivalParallel(ctx context.Context, recordMap map[string]*DemographyRecord, wg *sync.WaitGroup) {
	defer wg.Done() // сообщаем, что горутина завершилась
	c.loadArrival(ctx, recordMap)
}

// loadDepartureParallel - загружает убывших (отдельная горутина)
func (c *DataCollector) loadDepartureParallel(ctx context.Context, recordMap map[string]*DemographyRecord, wg *sync.WaitGroup) {
	defer wg.Done() // сообщаем, что горутина завершилась
	c.loadDeparture(ctx, recordMap)
}

// AggregateMigrationWithNormalize агрегирует миграцию с нормализацией OKTMO
func AggregateMigrationWithNormalize(records *[]MigrationRecord) map[string]map[int]int {
	result := make(map[string]map[int]int)
	for _, r := range *records {
		normalizedOktmo := NormalizeOktmo(r.Oktmo)

		if result[normalizedOktmo] == nil {
			result[normalizedOktmo] = make(map[int]int)
		}
		result[normalizedOktmo][r.Year] += r.Value
	}
	return result
}
