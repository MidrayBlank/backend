package parser

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
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
	// 1. Загружаем демографию (базовые данные)
	demRecords, err := c.loadDemography()
	if err != nil {
		return nil, fmt.Errorf("demography error: %w", err)
	}

	// Создаём карту для быстрого доступа
	recordMap := c.buildRecordMap(demRecords)

	// 2. Загружаем все дополнительные показатели
	c.loadLandArea(recordMap)
	c.loadHealthcare(recordMap)
	c.loadEducation(recordMap)
	c.loadHousing(recordMap)
	c.loadSalary(recordMap)
	c.loadAgeSex(recordMap)
	c.loadArrival(recordMap)
	c.loadDeparture(recordMap)

	// 3. Преобразуем в итоговый формат
	return c.convertToDBRecords(recordMap), nil
}

func (c *DataCollector) loadDemography() ([]DemographyRawRecord, error) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_mun_demography_123_v20240612/data_mun_demography_123_v20240612_csv.zip"

	file, err := DownloadAndUnzip(url)
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

		var births, deaths int
		if r.Births != 0 {
			births = r.Births
		}
		if r.Deaths != 0 {
			deaths = r.Deaths
		}

		recordMap[key] = &DemographyRecord{
			Code:             r.Oktmo,
			Year:             r.Year,
			PopulationAmount: r.Population,
			BirthAmount:      births,
			DeathAmount:      deaths,
		}
	}

	return recordMap
}

func (c *DataCollector) loadLandArea(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section6/data_Y48006001_112_v20250918.zip"

	file, err := DownloadAndUnzip(url)
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
				d.LandArea = value
			}
		}
	}
}

func (c *DataCollector) loadHealthcare(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section33/data_Y48018000_112_v20250918.zip"

	file, err := DownloadAndUnzip(url)
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
				val := int(value)
				d.MedicalFacilities = val
			}
		}
	}
}

func (c *DataCollector) loadEducation(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section34/data_Y48015001_112_v20250918.zip"

	file, err := DownloadAndUnzip(url)
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
				val := int(value)
				d.SchoolsCount = val
			}
		}
	}
}

func (c *DataCollector) loadHousing(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section37/data_Y48010001_112_v20250918.zip"

	file, err := DownloadAndUnzip(url)
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
				d.HousingCommissioned = value
			}
		}
	}
}

func (c *DataCollector) loadSalary(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section32/data_Y48423007_112_v20250918.zip"

	file, err := DownloadAndUnzip(url)
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
				d.AvgSalary = value
			}
		}
	}
}

func (c *DataCollector) loadAgeSex(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112014_112_v20250918.zip"

	file, err := DownloadAndUnzip(url)
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
				d.MaleAmount = values.Male
				d.FemaleAmount = values.Female
			}
		}
	}
}

func (c *DataCollector) loadArrival(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112021_112_v20250918.zip"

	fmt.Println("Loading arrival data...")
	files, err := DownloadAndUnzipAll(url)
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
				d.ArrivalAmount = value
			}
		}
	}
	fmt.Printf("  Arrival: %d records\n", len(allRecords))
}

func (c *DataCollector) loadDeparture(recordMap map[string]*DemographyRecord) {
	url := "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112022_112_v20250918.zip"

	fmt.Println("Loading departure data...")
	files, err := DownloadAndUnzipAll(url)
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
				d.DepartureAmount = value
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
