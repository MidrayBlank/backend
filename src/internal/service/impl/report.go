package impl

import (
	"backend/src/internal/domain"
	"backend/src/internal/service/abstract"

	connection "backend/src/internal/db/abstract"
	repository "backend/src/internal/repository/abstract"

	growthDomain "backend/src/pkg/report/sync/population_growth/domain"
	growthCalc "backend/src/pkg/report/sync/population_growth/report"

	leadersDomain "backend/src/pkg/report/sync/population_growth_leaders/domain"
	leadersCalc "backend/src/pkg/report/sync/population_growth_leaders/report"

	densityDomain "backend/src/pkg/report/sync/population_density/domain"
	densityCalc "backend/src/pkg/report/sync/population_density/report"

	indicatorsDomain "backend/src/pkg/report/sync/key_demographic_indicators/domain"
	indicatorsCalc "backend/src/pkg/report/sync/key_demographic_indicators/report"
)

type ReportService struct {
	conn        connection.IDBConnection
	rosstatRepo repository.IRosstatRepository
}

func NewReportService(conn connection.IDBConnection, rosstatRepo repository.IRosstatRepository) abstract.IReportService {
	return &ReportService{
		conn:        conn,
		rosstatRepo: rosstatRepo,
	}
}

func (r *ReportService) GetGrowthReport(codes []int, yearFrom, yearTo int) ([]domain.GrowthReport, error) {
	rosstatData, err := r.rosstatRepo.GetRosstatByCodes(r.conn, codes)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[int]*growthDomain.PopulationGrowthParams)

	for _, item := range rosstatData {
		if _, exists := dataMap[item.Code]; !exists {
			dataMap[item.Code] = &growthDomain.PopulationGrowthParams{
				OKTMO:    item.Code,
				YearFrom: yearFrom,
				YearTo:   yearTo,
			}
		}
		if item.Year == yearFrom {
			dataMap[item.Code].PopulationAmountFrom = getIntValue(item.PopulationAmount)
		}
		if item.Year == yearTo {
			dataMap[item.Code].PopulationAmountTo = getIntValue(item.PopulationAmount)
		}
	}

	compileData := make([]*growthDomain.PopulationGrowthParams, 0, len(dataMap))
	for _, v := range dataMap {
		if v.PopulationAmountFrom == 0 || v.PopulationAmountTo == 0 {
			continue
		}
		compileData = append(compileData, v)
	}

	compiler := growthCalc.NewPopulationGrowthReportCompiler()
	reports := compiler.Compile(compileData)

	result := make([]domain.GrowthReport, len(reports))
	for i, r := range reports {
		result[i] = domain.GrowthReport{
			Code:                r.OKTMO,
			YearFrom:            r.YearFrom,
			YearTo:              r.YearTo,
			PopulationGrowthPct: r.PopulationGrowthPct,
		}
	}

	return result, nil
}

func (s *ReportService) GetDensityReport(codes []int, year int) ([]domain.DensityReport, error) {

	rosstatData, err := s.rosstatRepo.GetRosstatByCodes(s.conn, codes)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[int]*densityDomain.PopulationDensityParams)

	for _, item := range rosstatData {
		if item.Year != year {
			continue
		}

		if _, exists := dataMap[item.Code]; !exists {
			dataMap[item.Code] = &densityDomain.PopulationDensityParams{
				OKTMO:            item.Code,
				Year:             year,
				PopulationAmount: getIntValue(item.PopulationAmount),
				LandArea:         float64(getIntValue(item.LandArea)),
			}
		}
	}

	compileData := make([]*densityDomain.PopulationDensityParams, 0, len(dataMap))
	for _, v := range dataMap {
		compileData = append(compileData, v)
	}

	calculator := densityCalc.NewPopulationDensityReportCompiler()
	reports := calculator.Compile(compileData)

	result := make([]domain.DensityReport, len(reports))
	for i, r := range reports {
		result[i] = domain.DensityReport{
			Code:    r.OKTMO,
			Density: r.Density,
		}
	}

	return result, nil
}

func (s *ReportService) GetGrowthLeadersReport(codes []int, yearFrom, yearTo, limit int) (*domain.GrowthLeadersResponse, error) {
	rosstatData, err := s.rosstatRepo.GetRosstatByCodes(s.conn, codes)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[int]*leadersDomain.GrowthLeadersParams)

	for _, item := range rosstatData {
		if _, exists := dataMap[item.Code]; !exists {
			dataMap[item.Code] = &leadersDomain.GrowthLeadersParams{
				OKTMO:    item.Code,
				YearFrom: yearFrom,
				YearTo:   yearTo,
				Limit:    limit,
			}
		}

		if item.Year == yearFrom {
			dataMap[item.Code].PopulationAmountFrom = getIntValue(item.PopulationAmount)
		}
		if item.Year == yearTo {
			dataMap[item.Code].PopulationAmountTo = getIntValue(item.PopulationAmount)
		}
	}

	compileData := make([]*leadersDomain.GrowthLeadersParams, 0, len(dataMap))
	for _, v := range dataMap {
		if v.PopulationAmountFrom == 0 || v.PopulationAmountTo == 0 {
			continue
		}
		compileData = append(compileData, v)
	}

	if len(compileData) == 0 {
		return &domain.GrowthLeadersResponse{
			TopGrowth:  []domain.GrowthLeadersReport{},
			TopDecline: []domain.GrowthLeadersReport{},
		}, nil
	}

	calculator := leadersCalc.NewGrowthLeadersReportCompiler()
	topGrowth, topDecline := calculator.Compile(compileData)

	growthResult := make([]domain.GrowthLeadersReport, len(topGrowth))
	for i, r := range topGrowth {
		growthResult[i] = domain.GrowthLeadersReport{
			Code:                r.OKTMO,
			YearFrom:            r.YearFrom,
			YearTo:              r.YearTo,
			PopulationGrowthPct: r.PopulationGrowthPct,
		}
	}

	declineResult := make([]domain.GrowthLeadersReport, len(topDecline))
	for i, r := range topDecline {
		declineResult[i] = domain.GrowthLeadersReport{
			Code:                r.OKTMO,
			YearFrom:            r.YearFrom,
			YearTo:              r.YearTo,
			PopulationGrowthPct: r.PopulationGrowthPct,
		}
	}

	return &domain.GrowthLeadersResponse{
		TopGrowth:  growthResult,
		TopDecline: declineResult,
	}, nil
}

func (s *ReportService) GetIndicatorReport(codes []int, year int) ([]domain.IndicatorsReport, error) {
	rosstatData, err := s.rosstatRepo.GetRosstatByCodes(s.conn, codes)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[int]*indicatorsDomain.DemographicParams)

	for _, item := range rosstatData {
		if item.Year != year {
			continue
		}

		if _, exists := dataMap[item.Code]; !exists {
			dataMap[item.Code] = &indicatorsDomain.DemographicParams{
				OKTMO:            item.Code,
				Year:             year,
				PopulationAmount: getIntValue(item.PopulationAmount),
				BirthAmount:      getIntValue(item.BirthAmount),
				DeathAmount:      getIntValue(item.DeathAmount),
				ArrivalAmount:    getIntValue(item.ArrivalAmount),
				DepartureAmount:  getIntValue(item.DepartureAmount),
			}
		}
	}

	compileData := make([]*indicatorsDomain.DemographicParams, 0, len(dataMap))
	for _, v := range dataMap {
		compileData = append(compileData, v)
	}

	if len(compileData) == 0 {
		return []domain.IndicatorsReport{}, nil
	}

	calculator := indicatorsCalc.NewDemographicReportCompiler()
	reports := calculator.Compile(compileData)

	result := make([]domain.IndicatorsReport, len(reports))
	for i, r := range reports {
		result[i] = domain.IndicatorsReport{
			Code:          r.OKTMO,
			Year:          r.Year,
			BirthRate:     r.BirthRate,
			DeathRate:     r.DeathRate,
			NaturalRate:   r.NaturalRate,
			MigrationRate: r.MigrationRate,
		}
	}

	return result, nil
}

func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
