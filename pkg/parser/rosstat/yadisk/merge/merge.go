package merge

import (
	"backend/pkg/parser/rosstat/yadisk/downloader"
	"backend/pkg/parser/rosstat/yadisk/model"
	"fmt"
)

// RecordMapManager управляет картой записей
type RecordMapManager struct {
	recordMap map[string]*model.DemographyRecord
}

// NewRecordMapManager создаёт новый менеджер
func NewRecordMapManager() *RecordMapManager {
	return &RecordMapManager{
		recordMap: make(map[string]*model.DemographyRecord),
	}
}

// GetRecordMap возвращает карту записей
func (r *RecordMapManager) GetRecordMap() map[string]*model.DemographyRecord {
	return r.recordMap
}

// GetOrCreate возвращает существующую запись или создаёт новую
func (r *RecordMapManager) GetOrCreate(oktmo string, year int) *model.DemographyRecord {
	key := fmt.Sprintf("%s_%d", oktmo, year)
	if _, ok := r.recordMap[key]; !ok {
		r.recordMap[key] = &model.DemographyRecord{
			Code: oktmo,
			Year: year,
		}
	}
	return r.recordMap[key]
}

// SetPopulation устанавливает численность населения
func (r *RecordMapManager) SetPopulation(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.PopulationAmount = value
}

// SetBirths устанавливает количество родившихся
func (r *RecordMapManager) SetBirths(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.BirthAmount = downloader.IntPtr(value)
}

// SetDeaths устанавливает количество умерших
func (r *RecordMapManager) SetDeaths(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.DeathAmount = downloader.IntPtr(value)
}

// SetLandArea устанавливает площадь
func (r *RecordMapManager) SetLandArea(oktmo string, year int, value float64) {
	record := r.GetOrCreate(oktmo, year)
	record.LandArea = downloader.Float64Ptr(value)
}

// SetMedicalFacilities устанавливает количество медучреждений
func (r *RecordMapManager) SetMedicalFacilities(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.MedicalFacilities = downloader.IntPtr(value)
}

// SetSchoolsCount устанавливает количество школ
func (r *RecordMapManager) SetSchoolsCount(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.SchoolsCount = downloader.IntPtr(value)
}

// SetHousingCommissioned устанавливает введённое жильё
func (r *RecordMapManager) SetHousingCommissioned(oktmo string, year int, value float64) {
	record := r.GetOrCreate(oktmo, year)
	record.HousingCommissioned = downloader.Float64Ptr(value)
}

// SetAvgSalary устанавливает среднюю зарплату
func (r *RecordMapManager) SetAvgSalary(oktmo string, year int, value float64) {
	record := r.GetOrCreate(oktmo, year)
	record.AvgSalary = downloader.Float64Ptr(value)
}

// SetAgeSex устанавливает данные по полу (суммарные)
func (r *RecordMapManager) SetAgeSex(oktmo string, year int, male, female int) {
	record := r.GetOrCreate(oktmo, year)
	record.MaleAmount = downloader.IntPtr(male)
	record.FemaleAmount = downloader.IntPtr(female)
}

// SetArrival устанавливает прибывших
func (r *RecordMapManager) SetArrival(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.ArrivalAmount = downloader.IntPtr(value)
}

// SetDeparture устанавливает убывших
func (r *RecordMapManager) SetDeparture(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.DepartureAmount = downloader.IntPtr(value)
}

// GetRecords возвращает слайс записей
func (r *RecordMapManager) GetRecords() []*model.DemographyRecord {
	result := make([]*model.DemographyRecord, 0, len(r.recordMap))
	for _, record := range r.recordMap {
		result = append(result, record)
	}
	return result
}
