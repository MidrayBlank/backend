package dto

type CreateReportRequest struct {
	Code int `json:"code" binding:"required"`
}

type ReportCodeParams struct {
	Code int `uri:"code" binding:"required"`
}

type ReportHashParams struct {
	Hash string `uri:"hash" binding:"required"`
}

type CreateReportResponse struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

type RequestStatusResponse struct {
	Hash   string `json:"hash"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ReportResponse struct {
	Report AIReportContent `json:"report"`
}

type AIReportContent struct {
	Summary         string   `json:"summary"`
	Forecast        Forecast `json:"forecast"`
	Trends          []string `json:"trends"`
	Recommendations []string `json:"recommendations"`
}

type Forecast struct {
	Scenario        string `json:"scenario"`
	Population      string `json:"population"`
	NaturalGrowth   string `json:"natural_growth"`
	MigrationGrowth string `json:"migration_growth"`
}
