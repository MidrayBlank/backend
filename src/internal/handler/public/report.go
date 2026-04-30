package public

import (
	context "backend/src/internal/context/abstract"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
	"strconv"
	"strings"
)

func GetGrowthReportHandler(ctx context.HandlerContext, params dto.GrowthReportParams, reportService service.IReportService) ([]dto.GrowthReportResponse, error) {
	codes := parseCodes(params.Codes)
	reports, err := reportService.GetGrowthReport(codes, params.YearFrom, params.YearTo)
	if err != nil {
		return nil, err
	}

	response := make([]dto.GrowthReportResponse, len(reports))
	for i, r := range reports {
		response[i] = dto.GrowthReportResponse{
			Code:                r.Code,
			YearFrom:            r.YearFrom,
			YearTo:              r.YearTo,
			PopulationGrowthPct: r.PopulationGrowthPct,
		}
	}

	return response, nil
}

func GetGrowthLeadersReportHandler(ctx context.HandlerContext, params dto.GrowthLeadersReportParams, reportService service.IReportService) (dto.GrowthLeadersResponse, error) {

	codes := parseCodes(params.Codes)
	limit := params.Limit

	leaders, err := reportService.GetGrowthLeadersReport(codes, params.YearFrom, params.YearTo, limit)
	if err != nil {
		return dto.GrowthLeadersResponse{}, err
	}

	response := dto.GrowthLeadersResponse{
		TopGrowth:  make([]dto.GrowthLeadersReportResponse, len(leaders.TopGrowth)),
		TopDecline: make([]dto.GrowthLeadersReportResponse, len(leaders.TopDecline)),
	}

	for i, r := range leaders.TopGrowth {
		response.TopGrowth[i] = dto.GrowthLeadersReportResponse{
			Code:                r.Code,
			YearFrom:            r.YearFrom,
			YearTo:              r.YearTo,
			PopulationGrowthPct: r.PopulationGrowthPct,
		}
	}

	for i, r := range leaders.TopDecline {
		response.TopDecline[i] = dto.GrowthLeadersReportResponse{
			Code:                r.Code,
			YearFrom:            r.YearFrom,
			YearTo:              r.YearTo,
			PopulationGrowthPct: r.PopulationGrowthPct,
		}
	}

	return response, nil
}

func GetDensityReportHandler(ctx context.HandlerContext, params dto.DensityReportParams, reportService service.IReportService) ([]dto.DensityReportResponse, error) {
	codes := parseCodes(params.Codes)
	reports, err := reportService.GetDensityReport(codes, params.Year)
	if err != nil {
		return nil, err
	}

	response := make([]dto.DensityReportResponse, len(reports))
	for i, r := range reports {
		response[i] = dto.DensityReportResponse{
			Code:    r.Code,
			Density: r.Density,
		}
	}

	return response, nil
}

func GetIndicatorsReportHandler(ctx context.HandlerContext, params dto.IndicatorsReportParams, reportService service.IReportService) ([]dto.IndicatorsReportResponse, error) {
	codes := parseCodes(params.Codes)
	reports, err := reportService.GetIndicatorReport(codes, params.Year)
	if err != nil {
		return nil, err
	}

	response := make([]dto.IndicatorsReportResponse, len(reports))
	for i, r := range reports {
		response[i] = dto.IndicatorsReportResponse{
			Code:          r.Code,
			Year:          r.Year,
			BirthRate:     r.BirthRate,
			DeathRate:     r.DeathRate,
			NaturalRate:   r.NaturalRate,
			MigrationRate: r.MigrationRate,
		}
	}

	return response, nil
}

func parseCodes(codesStr string) []int {
	var codes []int
	for _, s := range strings.Split(codesStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if code, err := strconv.Atoi(s); err == nil {
			codes = append(codes, code)
		}
	}
	return codes
}
