package storage

import (
	"strings"
	"sync"

	"backend/src/pkg/parser/rosstat/domain"
	"backend/src/pkg/parser/rosstat/rosstat/dao"
)

type storageKey struct {
	Code int
	Year int
}

type Storage struct {
	mapper       map[storageKey]*domain.RosstatParsed
	nameToCode   map[string]int
	subjectCodes [][]int
	m            *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		mapper:       make(map[storageKey]*domain.RosstatParsed),
		nameToCode:   make(map[string]int),
		subjectCodes: make([][]int, 100),
		m:            &sync.Mutex{},
	}
}

func (s *Storage) getRosstatParsed(code int, year int) *domain.RosstatParsed {
	key := storageKey{Code: code, Year: year}

	rosstatParsed, ok := s.mapper[key]
	if !ok {
		rosstatParsed = &domain.RosstatParsed{
			Code:       code,
			Year:       year,
			ParentCode: s.getParentCode(code),
		}
		s.mapper[key] = rosstatParsed
	}

	return rosstatParsed
}

func (s *Storage) SetCode(subjectCode int, extractedSlice []*dao.CodeExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		name := extracted.Name

		if idx := strings.Index(name, "("); idx != -1 {
			name = strings.TrimSpace(name[:idx])
		}

		s.nameToCode[name] = extracted.Code
		s.subjectCodes[subjectCode] = append(s.subjectCodes[subjectCode], extracted.Code)
	}
}

func (s *Storage) GetCode(name string) int {
	s.m.Lock()
	defer s.m.Unlock()

	if idx := strings.Index(name, "("); idx != -1 {
		name = strings.TrimSpace(name[:idx])
	}

	return s.nameToCode[name]
}

func (s *Storage) GetSubjectCodes(subjectCode int) []int {
	s.m.Lock()
	defer s.m.Unlock()

	return s.subjectCodes[subjectCode]
}

func (s *Storage) SetPopulation(extractedSlice []*dao.PopulationExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)

		rosstatParsed.Population = makeIntPtr(extracted.Population)
	}
}

func (s *Storage) SetBirth(extractedSlice []*dao.BirthExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Birth = makeIntPtr(extracted.Birth)
	}
}

func (s *Storage) SetDeath(extractedSlice []*dao.DeathExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Death = makeIntPtr(extracted.Death)
	}
}

func (s *Storage) SetArrival(extractedSlice []*dao.ArrivalExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Arrival = makeIntPtr(extracted.Arrival)
	}
}

func (s *Storage) SetDeparture(extractedSlice []*dao.DepartureExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Departure = makeIntPtr(extracted.Departure)
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
		rosstatParsed.Male = makeIntPtr(maleAmount)
		rosstatParsed.Female = makeIntPtr(femaleAmount)
	}
}

func (s *Storage) SetLandArea(extractedSlice []*dao.LandAreaExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.LandArea = makeIntPtr(extracted.LandArea)
	}
}

func (s *Storage) SetAverageSalary(extractedSlice []*dao.AverageSalaryExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.AverageSalary = makeFloatPtr(extracted.AverageSalary)
	}
}

func (s *Storage) SetMedicialFacilities(extractedSlice []*dao.MedicialFacilitiesExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.MedicialFacilities = makeIntPtr(extracted.MedicialFacilities)
	}
}

func (s *Storage) SetSchools(extractedSlice []*dao.SchoolsExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.Schools = makeIntPtr(extracted.Schools)
	}
}

func (s *Storage) SetHousingCommissioned(extractedSlice []*dao.HousingCommissionedExtracted) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, extracted := range extractedSlice {
		rosstatParsed := s.getRosstatParsed(extracted.Code, extracted.Year)
		rosstatParsed.HousingCommissioned = makeIntPtr(extracted.HousingCommissioned)
	}
}

func (s *Storage) Result() []*domain.RosstatParsed {
	result := make([]*domain.RosstatParsed, 0, len(s.mapper))
	for _, record := range s.mapper {
		result = append(result, record)
	}
	return result
}

func (s *Storage) getParentCode(code int) int {
	// if code%1000 != 0 {
	// 	return code / 1000 * 1000
	// }

	extendedSubjectCode := code / 100000
	subjectCode := code / 1000000

	switch extendedSubjectCode {
	case 118:
		subjectCode = extendedSubjectCode
	case 718:
		subjectCode = extendedSubjectCode
	case 719:
		subjectCode = extendedSubjectCode
	}

	return subjectCode
}
