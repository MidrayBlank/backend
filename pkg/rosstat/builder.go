package rosstat

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// joinInts преобразует срез целых чисел в строку с разделителем sep
func joinInts(nums []int, sep string) string {
	if len(nums) == 0 {
		return ""
	}
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, sep)
}

// buildQryString собирает строку для параметра Qry
func buildQryString(qry QueryDimensions) string {
	var parts []string
	if len(qry.Pokazateli) > 0 {
		parts = append(parts, "Pokazateli:"+joinInts(qry.Pokazateli, ",")+";")
	}
	if len(qry.Munr) > 0 {
		parts = append(parts, "munr:"+joinInts(qry.Munr, ",")+";")
	}
	if len(qry.Tippos) > 0 {
		parts = append(parts, "tippos:"+joinInts(qry.Tippos, ",")+";")
	}
	if len(qry.Oktmo) > 0 {
		parts = append(parts, "oktmo:"+joinInts(qry.Oktmo, ",")+";")
	}
	parts = append(parts, "vozr:"+strconv.Itoa(qry.Vozr)+";")
	if len(qry.Grup_2) > 0 {
		parts = append(parts, "grup_2:"+joinInts(qry.Grup_2, ",")+";")
	}
	if len(qry.God) > 0 {
		parts = append(parts, "god:"+joinInts(qry.God, ",")+";")
	}
	parts = append(parts, "period:"+strconv.Itoa(qry.Period)+";")
	if len(qry.Mest) > 0 {
		parts = append(parts, "mest:"+joinInts(qry.Mest, ",")+";")
	}
	return strings.Join(parts, "")
}

// buildQryGmString собирает строку для параметра QryGm
func buildQryGmString(qryGm QueryGm) string {
	var parts []string
	parts = append(parts, "vozr_z:"+strconv.Itoa(qryGm.Vozr_z)+";")
	parts = append(parts, "period_z:"+strconv.Itoa(qryGm.Period_z)+";")
	parts = append(parts, "god_s:"+strconv.Itoa(qryGm.God_s)+";")
	parts = append(parts, "munr_b:"+strconv.Itoa(qryGm.Munr_b)+";")
	parts = append(parts, "oktmo_b:"+strconv.Itoa(qryGm.Oktmo_b)+";")
	parts = append(parts, "grup_2_b:"+strconv.Itoa(qryGm.Grup_2_b)+";")
	parts = append(parts, "Pokazateli_b:"+strconv.Itoa(qryGm.Pokazateli_b)+";")
	parts = append(parts, "tippos_b:"+strconv.Itoa(qryGm.Tippos_b)+";")
	parts = append(parts, "mest_b:"+strconv.Itoa(qryGm.Mest_b)+";")
	return strings.Join(parts, "")
}

// BuildRequestBody преобразует RequestParams в строку тела POST-запроса
func BuildRequestBody(params RequestParams) string {
	qryRaw := buildQryString(params.Qry)
	qryEncoded := url.QueryEscape(qryRaw)
	qryGmRaw := buildQryGmString(params.QryGm)
	qryGmEncoded := url.QueryEscape(qryGmRaw)
	yearsListRaw := joinInts(params.YearsList, ";") + ";"
	yearsListEncoded := url.QueryEscape(yearsListRaw)
	body := fmt.Sprintf(
		"DiagSz=%s&Format=%s&YearTo=%d&YearFrom=%d&Qry=%s&QryGm=%s&QryFootNotes=%%3B&YearsList=%s&tbl=%s",
		params.DiagSz,
		params.Format,
		params.YearTo,
		params.YearFrom,
		qryEncoded,
		qryGmEncoded,
		yearsListEncoded,
		"%CF%EE%EA%E0%E7%E0%F2%FC+%F2%E0%E1%EB%E8%F6%F3",
	)
	return body
}
