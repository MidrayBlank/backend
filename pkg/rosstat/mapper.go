package rosstat

// NewRequestParams создаёт полностью заполненную структуру RequestParams и возвращает ее для отправки запроса

func NewRequestParams(territoryCodes []int, yearFrom, yearTo int) RequestParams {
	years := generateYears(yearFrom, yearTo)
	return RequestParams{
		Format:   "CSV",
		DiagSz:   "800x600",
		Tbl:      "Показать таблицу",
		YearFrom: yearFrom,
		YearTo:   yearTo,
		Qry: QueryDimensions{
			Pokazateli: []int{8112027, 8112014},
			Munr:       territoryCodes,
			Oktmo:      territoryCodes,
			Tippos:     []int{10, 7, 1, 4, 20},
			Vozr:       151,
			Grup_2:     []int{1, 2, 3},
			God:        years,
			Period:     208,
			Mest:       []int{11, 12, 13},
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
		QryFootNotes: "",
		YearsList:    years,
	}
}

func generateYears(from, to int) []int {
	years := make([]int, 0, to-from+1)
	for i := from; i <= to; i++ {
		years = append(years, i)
	}
	return years
}
