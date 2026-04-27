package yadisk_parse

import (
	"context"
	"fmt"
	"log"

	"backend/src/internal/async/semaphore"
	"backend/src/internal/async/status"
	"backend/src/internal/async/worker/repository"
	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"

	parser "backend/src/pkg/parser/rosstat/domain"
	"backend/src/pkg/parser/rosstat/yadisk"
)

const batchSize = 2000

func YadiskParseWorker(
	ctx context.Context,
	ch chan status.CompletionStatus,
	sem semaphore.Semaphore,
	conn abstract.IDBConnection,
	repositories repository.WorkerRepositories,
	requestId int,
	param int,
) {
	defer sem.Release()
	parser := yadisk.NewYadiskRosstatParser()

	log.Printf("Parsing...\n")
	rosstatParsedSlice, err := parser.Parse(ctx)
	if err != nil {
		log.Printf("Request %d failed with error: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	log.Printf("Upsert GEO...\n")
	if err := batchUpsertGeo(ctx, conn, repositories, rosstatParsedSlice); err != nil {
		log.Printf("Request %d failed in upsertGeo: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	log.Printf("Aggregate by subject...\n")
	subjectAggs, err := aggregateBySubject(rosstatParsedSlice)
	if err != nil {
		log.Printf("Request %d failed in aggregateBySubject: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	log.Printf("Aggregate RF...\n")
	rfAgg, err := aggregateRF(subjectAggs)
	if err != nil {
		log.Printf("Request %d failed in aggregateRF: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	log.Printf("Upsert Rosstat with age...\n")

	var allRosstat []*domain.Rosstat
	var allAge []*domain.RosstatByAge

	for _, p := range rosstatParsedSlice {
		rosstat := &domain.Rosstat{
			Code:                p.Code,
			Year:                p.Year,
			PopulationAmount:    p.Population,
			BirthAmount:         p.Birth,
			DeathAmount:         p.Death,
			ArrivalAmount:       p.Arrival,
			DepartureAmount:     p.Departure,
			MaleAmount:          p.Male,
			FemaleAmount:        p.Female,
			LandArea:            p.LandArea,
			AvgSalary:           p.AverageSalary,
			MedicalFacilities:   p.MedicialFacilities,
			SchoolsCount:        p.Schools,
			HousingCommissioned: p.HousingCommissioned,
		}
		allRosstat = append(allRosstat, rosstat)
	}

	subjectRosstat := make(map[int]*domain.Rosstat)
	for subjCode, agg := range subjectAggs {
		rosstat := &domain.Rosstat{
			Code:                subjCode,
			Year:                agg.Year,
			PopulationAmount:    agg.Population,
			BirthAmount:         agg.Birth,
			DeathAmount:         agg.Death,
			ArrivalAmount:       agg.Arrival,
			DepartureAmount:     agg.Departure,
			MaleAmount:          agg.Male,
			FemaleAmount:        agg.Female,
			LandArea:            agg.LandArea,
			AvgSalary:           agg.AvgSalary,
			MedicalFacilities:   agg.MedicalFacilities,
			SchoolsCount:        agg.Schools,
			HousingCommissioned: agg.HousingCommissioned,
		}
		allRosstat = append(allRosstat, rosstat)
		subjectRosstat[subjCode] = rosstat
	}

	rfRosstat := &domain.Rosstat{
		Code:                0,
		Year:                rfAgg.Year,
		PopulationAmount:    rfAgg.Population,
		BirthAmount:         rfAgg.Birth,
		DeathAmount:         rfAgg.Death,
		ArrivalAmount:       rfAgg.Arrival,
		DepartureAmount:     rfAgg.Departure,
		MaleAmount:          rfAgg.Male,
		FemaleAmount:        rfAgg.Female,
		LandArea:            rfAgg.LandArea,
		AvgSalary:           rfAgg.AvgSalary,
		MedicalFacilities:   rfAgg.MedicalFacilities,
		SchoolsCount:        rfAgg.Schools,
		HousingCommissioned: rfAgg.HousingCommissioned,
	}
	allRosstat = append(allRosstat, rfRosstat)

	if err := batchUpsertRosstat(ctx, conn, repositories, allRosstat); err != nil {
		log.Printf("Request %d failed during batch upsert Rosstat: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	codesSet := make(map[int]bool)
	for _, r := range allRosstat {
		codesSet[r.Code] = true
	}
	codes := make([]int, 0, len(codesSet))
	for c := range codesSet {
		codes = append(codes, c)
	}
	rosstatList, err := repositories.RosstatRepository.GetRosstatByCodes(conn, codes)
	if err != nil {
		log.Printf("Request %d failed to fetch Rosstat IDs: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}
	rosstatIDMap := make(map[[2]int]int)
	for _, r := range rosstatList {
		key := [2]int{r.Code, r.Year}
		rosstatIDMap[key] = r.ID
	}

	for _, p := range rosstatParsedSlice {
		key := [2]int{p.Code, p.Year}
		rosstatID, ok := rosstatIDMap[key]
		if !ok {
			log.Printf("Rosstat ID not found for code=%d year=%d\n", p.Code, p.Year)
			continue
		}
		for _, a := range p.ByAge {
			allAge = append(allAge, &domain.RosstatByAge{
				RosstatID:    rosstatID,
				Age:          a.Age,
				MaleAmount:   a.MaleAmount,
				FemaleAmount: a.FemaleAmount,
			})
		}
	}
	for subjCode, agg := range subjectAggs {
		key := [2]int{subjCode, agg.Year}
		rosstatID, ok := rosstatIDMap[key]
		if !ok {
			log.Printf("Rosstat ID not found for subject code=%d year=%d\n", subjCode, agg.Year)
			continue
		}
		for _, a := range agg.ByAge {
			allAge = append(allAge, &domain.RosstatByAge{
				RosstatID:    rosstatID,
				Age:          a.Age,
				MaleAmount:   a.MaleAmount,
				FemaleAmount: a.FemaleAmount,
			})
		}
	}
	{
		key := [2]int{0, rfAgg.Year}
		rosstatID, ok := rosstatIDMap[key]
		if !ok {
			log.Printf("Rosstat ID not found for RF year=%d\n", rfAgg.Year)
		} else {
			for _, a := range rfAgg.ByAge {
				allAge = append(allAge, &domain.RosstatByAge{
					RosstatID:    rosstatID,
					Age:          a.Age,
					MaleAmount:   a.MaleAmount,
					FemaleAmount: a.FemaleAmount,
				})
			}
		}
	}

	if err := batchUpsertRosstatAge(ctx, conn, repositories, allAge); err != nil {
		log.Printf("Request %d failed during batch upsert RosstatByAge: %v\n", requestId, err)
		ch <- status.CompletionStatus{RequestId: requestId, Err: err}
		return
	}

	ch <- status.CompletionStatus{RequestId: requestId, Err: nil}
}

func batchUpsertGeo(
	ctx context.Context,
	conn abstract.IDBConnection,
	repos repository.WorkerRepositories,
	parsedSlice parser.RosstatParsedSlice,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	seen := make(map[int]bool)
	var allGeo []*domain.Geo

	for _, p := range parsedSlice {
		if seen[p.Code] {
			continue
		}
		seen[p.Code] = true
		level, parentCode := municipalLevelAndParent(p.Code, p.SubjectCode)
		geo := &domain.Geo{
			Code:       p.Code,
			ParentCode: parentCode,
			Name:       "",
			Level:      level,
		}
		allGeo = append(allGeo, geo)
	}

	for i := 0; i < len(allGeo); i += batchSize {
		end := i + batchSize
		if end > len(allGeo) {
			end = len(allGeo)
		}
		batch := allGeo[i:end]
		if err := repos.GeoRepository.UpsertBatch(conn, batch); err != nil {
			return fmt.Errorf("geo upsert batch: %w", err)
		}
	}
	return nil
}

func aggregateBySubject(parsed parser.RosstatParsedSlice) (map[int]*aggregatedData, error) {
	byKey := make(map[int]*aggregatedData)

	for _, p := range parsed {
		subj := p.SubjectCode
		agg, ok := byKey[subj]
		if !ok {
			agg = &aggregatedData{
				Year:    p.Year,
				ageSums: make(map[int]*ageSum),
			}
			byKey[subj] = agg
		}

		addInt(&agg.Population, p.Population)
		addInt(&agg.Birth, p.Birth)
		addInt(&agg.Death, p.Death)
		addInt(&agg.Arrival, p.Arrival)
		addInt(&agg.Departure, p.Departure)
		addInt(&agg.Male, p.Male)
		addInt(&agg.Female, p.Female)
		addInt(&agg.LandArea, nil)
		addInt(&agg.MedicalFacilities, p.MedicialFacilities)
		addInt(&agg.Schools, p.Schools)
		addInt(&agg.HousingCommissioned, p.HousingCommissioned)

		if p.AverageSalary != nil && p.Population != nil {
			agg.sumSalaryWeighted += *p.AverageSalary * float64(*p.Population)
			agg.totalPopForSalary += *p.Population
		}

		for _, a := range p.ByAge {
			sum := agg.ageSums[a.Age]
			if sum == nil {
				sum = &ageSum{}
				agg.ageSums[a.Age] = sum
			}
			sum.male += a.MaleAmount
			sum.female += a.FemaleAmount
		}
	}

	for _, agg := range byKey {
		if agg.totalPopForSalary > 0 {
			avg := agg.sumSalaryWeighted / float64(agg.totalPopForSalary)
			agg.AvgSalary = &avg
		}
		agg.ByAge = sortedAgeSlice(agg.ageSums)
		agg.ageSums = nil
	}

	return byKey, nil
}

func aggregateRF(subjectAggs map[int]*aggregatedData) (*aggregatedData, error) {
	rf := &aggregatedData{
		ageSums: make(map[int]*ageSum),
	}

	for _, agg := range subjectAggs {
		if rf.Year == 0 {
			rf.Year = agg.Year
		}

		addInt(&rf.Population, agg.Population)
		addInt(&rf.Birth, agg.Birth)
		addInt(&rf.Death, agg.Death)
		addInt(&rf.Arrival, agg.Arrival)
		addInt(&rf.Departure, agg.Departure)
		addInt(&rf.Male, agg.Male)
		addInt(&rf.Female, agg.Female)
		addInt(&rf.LandArea, nil)
		addInt(&rf.MedicalFacilities, agg.MedicalFacilities)
		addInt(&rf.Schools, agg.Schools)
		addInt(&rf.HousingCommissioned, agg.HousingCommissioned)

		if agg.AvgSalary != nil && agg.Population != nil {
			rf.sumSalaryWeighted += *agg.AvgSalary * float64(*agg.Population)
			rf.totalPopForSalary += *agg.Population
		}

		for _, a := range agg.ByAge {
			sum := rf.ageSums[a.Age]
			if sum == nil {
				sum = &ageSum{}
				rf.ageSums[a.Age] = sum
			}
			sum.male += a.MaleAmount
			sum.female += a.FemaleAmount
		}
	}

	if rf.totalPopForSalary > 0 {
		avg := rf.sumSalaryWeighted / float64(rf.totalPopForSalary)
		rf.AvgSalary = &avg
	}
	rf.ByAge = sortedAgeSlice(rf.ageSums)
	rf.ageSums = nil
	return rf, nil
}

func batchUpsertRosstat(
	ctx context.Context,
	conn abstract.IDBConnection,
	repos repository.WorkerRepositories,
	rosstatList []*domain.Rosstat,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for i := 0; i < len(rosstatList); i += batchSize {
		end := i + batchSize
		if end > len(rosstatList) {
			end = len(rosstatList)
		}
		batch := rosstatList[i:end]
		if err := repos.RosstatRepository.UpsertBatch(conn, batch); err != nil {
			return fmt.Errorf("rosstat upsert batch: %w", err)
		}
	}
	return nil
}

func batchUpsertRosstatAge(
	ctx context.Context,
	conn abstract.IDBConnection,
	repos repository.WorkerRepositories,
	ageList []*domain.RosstatByAge,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for i := 0; i < len(ageList); i += batchSize {
		end := i + batchSize
		if end > len(ageList) {
			end = len(ageList)
		}
		batch := ageList[i:end]
		if err := repos.RosstatAgeRepository.UpsertBatch(conn, batch); err != nil {
			return fmt.Errorf("rosstat age upsert batch: %w", err)
		}
	}
	return nil
}

func municipalLevelAndParent(code int, subjectCode int) (int, *int) {
	var level, parentCode int

	switch code % 1000 {
	case 0:
		level = 2
		parentCode = subjectCode
	default:
		level = 3
		parentCode = (code / 1000) * 1000
	}

	return level, &parentCode
}

type aggregatedData struct {
	Year                int
	Population          *int
	Birth               *int
	Death               *int
	Arrival             *int
	Departure           *int
	Male                *int
	Female              *int
	LandArea            *int
	AvgSalary           *float64
	MedicalFacilities   *int
	Schools             *int
	HousingCommissioned *int
	ByAge               []*parser.RosstatAgeParsed

	sumSalaryWeighted float64
	totalPopForSalary int
	ageSums           map[int]*ageSum
}

type ageSum struct {
	male   int
	female int
}

func addInt(dst **int, src *int) {
	if src == nil {
		return
	}
	if *dst == nil {
		val := 0
		*dst = &val
	}
	**dst += *src
}

func sortedAgeSlice(m map[int]*ageSum) []*parser.RosstatAgeParsed {
	if len(m) == 0 {
		return nil
	}
	ages := make([]int, 0, len(m))
	for a := range m {
		ages = append(ages, a)
	}
	for i := 0; i < len(ages); i++ {
		for j := i + 1; j < len(ages); j++ {
			if ages[j] < ages[i] {
				ages[i], ages[j] = ages[j], ages[i]
			}
		}
	}

	result := make([]*parser.RosstatAgeParsed, len(ages))
	for i, age := range ages {
		sum := m[age]
		result[i] = &parser.RosstatAgeParsed{
			Age:          age,
			MaleAmount:   sum.male,
			FemaleAmount: sum.female,
		}
	}
	return result
}
