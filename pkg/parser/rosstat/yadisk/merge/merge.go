package merge

import (
	"backend/pkg/parser/rosstat/yadisk/downloader"
	"backend/pkg/parser/rosstat/yadisk/model"
	"fmt"
)

type RecordMapManager struct {
	recordMap map[string]*model.RosstatParsed
}

func NewRecordMapManager() *RecordMapManager {
	return &RecordMapManager{
		recordMap: make(map[string]*model.RosstatParsed),
	}
}

func (r *RecordMapManager) GetRecordMap() map[string]*model.RosstatParsed {
	return r.recordMap
}

func (r *RecordMapManager) GetOrCreate(oktmo string, year int) *model.RosstatParsed {
	key := fmt.Sprintf("%s_%d", oktmo, year)
	if _, ok := r.recordMap[key]; !ok {
		r.recordMap[key] = &model.RosstatParsed{
			Code: oktmo,
			Year: year,
		}
	}
	return r.recordMap[key]
}

func (r *RecordMapManager) SetPopulation(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.PopulationAmount = value
}

func (r *RecordMapManager) SetBirths(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.BirthAmount = downloader.IntPtr(value)
}

func (r *RecordMapManager) SetDeaths(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.DeathAmount = downloader.IntPtr(value)
}

func (r *RecordMapManager) SetLandArea(oktmo string, year int, value float64) {
	record := r.GetOrCreate(oktmo, year)
	record.LandArea = downloader.Float64Ptr(value)
}

func (r *RecordMapManager) SetMedicalFacilities(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.MedicalFacilities = downloader.IntPtr(value)
}

func (r *RecordMapManager) SetSchoolsCount(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.SchoolsCount = downloader.IntPtr(value)
}

func (r *RecordMapManager) SetHousingCommissioned(oktmo string, year int, value float64) {
	record := r.GetOrCreate(oktmo, year)
	record.HousingCommissioned = downloader.Float64Ptr(value)
}

func (r *RecordMapManager) SetAvgSalary(oktmo string, year int, value float64) {
	record := r.GetOrCreate(oktmo, year)
	record.AvgSalary = downloader.Float64Ptr(value)
}

func (r *RecordMapManager) SetAgeSex(oktmo string, year int, male, female int) {
	record := r.GetOrCreate(oktmo, year)
	record.MaleAmount = downloader.IntPtr(male)
	record.FemaleAmount = downloader.IntPtr(female)
}

func (r *RecordMapManager) SetArrival(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.ArrivalAmount = downloader.IntPtr(value)
}

func (r *RecordMapManager) SetDeparture(oktmo string, year int, value int) {
	record := r.GetOrCreate(oktmo, year)
	record.DepartureAmount = downloader.IntPtr(value)
}

func (r *RecordMapManager) GetRecords() []*model.RosstatParsed {
	result := make([]*model.RosstatParsed, 0, len(r.recordMap))
	for _, record := range r.recordMap {
		result = append(result, record)
	}
	return result
}
