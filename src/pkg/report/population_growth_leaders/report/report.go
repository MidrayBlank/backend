package report

import (
	"backend/src/pkg/report/population_growth_leaders/domain"
	"backend/src/pkg/report/population_growth_leaders/dto"
	"sort"
)

type GrowthLeadersCalculator struct{}

func NewGrowthLeadersCalculator() *GrowthLeadersCalculator {
	return &GrowthLeadersCalculator{}
}

func (c *GrowthLeadersCalculator) CompileReport(data []*domain.GrowthLeadersParams) (topGrowth []*domain.GrowthLeadersReport, topDecline []*domain.GrowthLeadersReport) {
	if len(data) == 0 {
		return []*domain.GrowthLeadersReport{}, []*domain.GrowthLeadersReport{}
	}

	items := make([]dto.GrowthItem, 0, len(data))
	for _, d := range data {
		if d.PopulationAmountFrom == 0 {
			continue
		}
		growthPct := float32(d.PopulationAmountTo-d.PopulationAmountFrom) / float32(d.PopulationAmountFrom) * 100
		items = append(items, dto.GrowthItem{
			OKTMO:     d.OKTMO,
			YearFrom:  d.YearFrom,
			YearTo:    d.YearTo,
			GrowthPct: growthPct,
		})
	}
	if len(items) == 0 {
		return []*domain.GrowthLeadersReport{}, []*domain.GrowthLeadersReport{}
	}

	limit := data[0].Limit

	topGrowth = c.getTopByGrowth(items, limit, false)

	topDecline = c.getTopByGrowth(items, limit, true)

	return topGrowth, topDecline
}

func (c *GrowthLeadersCalculator) getTopByGrowth(items []dto.GrowthItem, limit int, flag bool) []*domain.GrowthLeadersReport {
	sorted := make([]dto.GrowthItem, len(items))
	copy(sorted, items)

	sort.Slice(sorted, func(i, j int) bool {
		if flag {
			return sorted[i].GrowthPct < sorted[j].GrowthPct
		}
		return sorted[i].GrowthPct > sorted[j].GrowthPct
	})

	resultLimit := limit
	if len(sorted) < limit {
		resultLimit = len(sorted)
	}

	report := make([]*domain.GrowthLeadersReport, resultLimit)
	for i := 0; i < resultLimit; i++ {
		report[i] = &domain.GrowthLeadersReport{
			OKTMO:               sorted[i].OKTMO,
			YearFrom:            sorted[i].YearFrom,
			YearTo:              sorted[i].YearTo,
			PopulationGrowthPct: sorted[i].GrowthPct,
		}
	}
	return report
}
