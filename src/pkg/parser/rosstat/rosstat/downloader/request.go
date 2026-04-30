package downloader

import (
	"backend/src/pkg/parser/rosstat/rosstat/config"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type QueryDimensions struct {
	Pokazateli []int
	Munr       []int
	Tippos     []int
	Oktmo      []int
	Vozr       int
	Grup_2     []int
	God        []int
	Period     int
	Mest       []int
}

type QueryGm struct {
	Vozr_z       int
	Period_z     int
	God_s        int
	Munr_b       int
	Oktmo_b      int
	Grup_2_b     int
	Pokazateli_b int
	Tippos_b     int
	Mest_b       int
}

type RequestParameters struct {
	Format    string
	YearFrom  int
	YearTo    int
	Qry       QueryDimensions
	QryGm     QueryGm
	YearsList []int
}

func NewRequestParameters(indicator int, munr []int, oktmo []int, yearTo int) *RequestParameters {
	config := config.NewConfig()
	years := generateYears(2000, yearTo)
	return &RequestParameters{
		Format:   "CSV",
		YearFrom: 2000,
		YearTo:   yearTo,
		Qry: QueryDimensions{
			Pokazateli: []int{indicator},
			Munr:       munr,
			Tippos:     []int{10, 7, 1, 4, 20},
			Oktmo:      oktmo,
			Vozr:       151,
			Grup_2:     []int{1, 2, 3},
			God:        years,
			Period:     config.GetPeriod(indicator),
			Mest:       []int{11},
		},
		QryGm: QueryGm{
			Vozr_z:       1,
			Period_z:     2,
			God_s:        1,
			Munr_b:       1,
			Oktmo_b:      2,
			Grup_2_b:     3,
			Pokazateli_b: 4,
			Tippos_b:     5,
			Mest_b:       6,
		},
		YearsList: years,
	}
}

func (p *RequestParameters) buildRequestBodyStr() string {
	qryRaw := fmt.Sprintf(
		"Pokazateli:%s;munr:%s;tippos:%s;oktmo:%s;vozr:%d;grup_2:%s;god:%s;period:%d;mest:%s;",
		joinInts(p.Qry.Pokazateli, ","),
		joinInts(p.Qry.Munr, ","),
		joinInts(p.Qry.Tippos, ","),
		joinInts(p.Qry.Oktmo, ","),
		p.Qry.Vozr,
		joinInts(p.Qry.Grup_2, ","),
		joinInts(p.Qry.God, ","),
		p.Qry.Period,
		joinInts(p.Qry.Mest, ","),
	)
	qryEncoded := url.QueryEscape(qryRaw)

	qryGmRaw := fmt.Sprintf(
		"vozr_z:%d;period_z:%d;god_s:%d;munr_b:%d;oktmo_b:%d;grup_2_b:%d;Pokazateli_b:%d;tippos_b:%d;mest_b:%d;",
		p.QryGm.Vozr_z, p.QryGm.Period_z, p.QryGm.God_s,
		p.QryGm.Munr_b, p.QryGm.Oktmo_b, p.QryGm.Grup_2_b,
		p.QryGm.Pokazateli_b, p.QryGm.Tippos_b, p.QryGm.Mest_b,
	)
	qryGmEncoded := url.QueryEscape(qryGmRaw)

	yearsRaw := joinInts(p.YearsList, ";") + ";"
	yearsEncoded := url.QueryEscape(yearsRaw)

	body := fmt.Sprintf(
		"DiagSz=800x600&Format=%s&YearTo=%d&YearFrom=%d&Qry=%s&QryGm=%s&QryFootNotes=%%3B&YearsList=%s&tbl=%s",
		p.Format,
		p.YearTo,
		p.YearFrom,
		qryEncoded,
		qryGmEncoded,
		yearsEncoded,
		"%CF%EE%EA%E0%E7%E0%F2%FC+%F2%E0%E1%EB%E8%F6%F3",
	)
	return body
}

func joinInts(nums []int, sep string) string {
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, sep)
}

func generateYears(yearFrom int, yearTo int) []int {
	years := make([]int, 0, yearTo-yearFrom+1)
	for y := yearFrom; y <= yearTo; y++ {
		years = append(years, y)
	}
	return years
}
