package storage

import (
	"sync"

	"backend/src/pkg/parser/rosstat/domain"
	"backend/src/pkg/parser/rosstat/yadisk/dao"
)

type storageKey struct {
	Code int
	Year int
}

type Storage struct {
	mapper map[storageKey]*domain.RosstatParsed
	m      *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		mapper: make(map[storageKey]*domain.RosstatParsed),
		m:      &sync.Mutex{},
	}
}

func (s *Storage) getRosstatParsed(code int, year int) *domain.RosstatParsed {
	key := storageKey{Code: code, Year: year}

	rosstatParsed, ok := s.mapper[key]
	if !ok {
		rosstatParsed = &domain.RosstatParsed{
			Code:        code,
			Year:        year,
			SubjectCode: s.getSubjectCode(code),
		}
		s.mapper[key] = rosstatParsed
	}

	return rosstatParsed
}

func (s *Storage) SetPopulation(extractedSlice []*dao.PopulationExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)

		if !rosstatParsed.RuralUsed && !rosstatParsed.UrbanUsed {
			rosstatParsed.Population = intPtr(extracted.Population)
			rosstatParsed.UrbanUsed = extracted.UrbanUsed
			rosstatParsed.RuralUsed = extracted.RuralUsed
		} else if !rosstatParsed.RuralUsed && rosstatParsed.UrbanUsed && extracted.RuralUsed && !extracted.UrbanUsed {
			*rosstatParsed.Population += extracted.Population
			rosstatParsed.RuralUsed = extracted.RuralUsed
		} else if rosstatParsed.RuralUsed && !rosstatParsed.UrbanUsed && !extracted.RuralUsed && extracted.UrbanUsed {
			*rosstatParsed.Population += extracted.Population
			rosstatParsed.UrbanUsed = extracted.UrbanUsed
		}
	}
}

func (s *Storage) SetBirth(extractedSlice []*dao.BirthExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Birth = intPtr(extracted.Birth)
	}
}

func (s *Storage) SetDeath(extractedSlice []*dao.DeathExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Death = intPtr(extracted.Death)
	}
}

func (s *Storage) SetArrival(extractedSlice []*dao.ArrivalExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Arrival = intPtr(extracted.Arrival)
	}
}

func (s *Storage) SetDeparture(extractedSlice []*dao.DepartureExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Departure = intPtr(extracted.Departure)
	}
}

func (s *Storage) SetMaleFemaleAge(extractedSlice []*dao.MaleFemaleAgeExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	groups := make(map[storageKey][]*dao.MaleFemaleAgeExtracted)
	for _, ex := range extractedSlice {
		key := storageKey{Code: ex.Code, Year: ex.Year}
		groups[key] = append(groups[key], ex)
	}

	for key, exGroup := range groups {
		rosstatParsed := s.getRosstatParsed(key.Code, key.Year)

		var maleAmount, femaleAmount int
		ageParsed := make([]*domain.RosstatAgeParsed, 0, len(exGroup))

		for _, ex := range exGroup {
			maleAmount += ex.MaleAmount
			femaleAmount += ex.FemaleAmount

			ageParsed = append(ageParsed, &domain.RosstatAgeParsed{
				Age:          ex.Age,
				MaleAmount:   ex.MaleAmount,
				FemaleAmount: ex.FemaleAmount,
			})
		}

		rosstatParsed.ByAge = ageParsed
		rosstatParsed.Male = intPtr(maleAmount)
		rosstatParsed.Female = intPtr(femaleAmount)
	}
}

func (s *Storage) SetLandArea(extractedSlice []*dao.LandAreaExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.LandArea = intPtr(extracted.LandArea)
	}
}

func (s *Storage) SetAverageSalary(extractedSlice []*dao.AverageSalaryExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.AverageSalary = float64Ptr(extracted.AverageSalary)
	}
}

func (s *Storage) SetMedicialFacilities(extractedSlice []*dao.MedicialFacilitiesExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.MedicialFacilities = intPtr(extracted.MedicialFacilities)
	}
}

func (s *Storage) SetSchools(extractedSlice []*dao.SchoolsExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Schools = intPtr(extracted.Schools)
	}
}

func (s *Storage) SetHousingCommissioned(extractedSlice []*dao.HousingCommissionedExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.HousingCommissioned = intPtr(extracted.HousingCommissioned)
	}
}

func (s *Storage) Result() []*domain.RosstatParsed {
	result := make([]*domain.RosstatParsed, 0, len(s.mapper))
	for _, record := range s.mapper {
		result = append(result, record)
	}
	return result
}

func (s *Storage) getSubjectCode(code int) int {
	extendedCode := code / 100000
	realCode := code / 1000000

	switch extendedCode {
	case 118:
		realCode = extendedCode
	case 718:
		realCode = extendedCode
	case 719:
		realCode = extendedCode
	}

	return realCode
}

func intPtr(v int) *int {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}
