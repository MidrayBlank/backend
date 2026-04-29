package public

import (
	request_type "backend/src/internal/async/request"
	context "backend/src/internal/context/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	service "backend/src/internal/service/abstract"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetReportAsyncHandler(ctx context.HandlerContext,
	aiReportAsyncService service.IAiReportAsyncService,
) dto.ReportResponse {
	codeStr := ctx.Get("code")

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return dto.ReportResponse{}
	}

	if code <= 0 {
		ctx.Status(http.StatusBadRequest)
		return dto.ReportResponse{}
	}

	report, err := aiReportAsyncService.GetReportByCode(code)
	if err != nil || report == nil {
		ctx.Status(http.StatusBadRequest)
		return dto.ReportResponse{}
	}

	return buildAsyncReportResponse(ctx, report)
}

func buildAsyncReportResponse(ctx context.HandlerContext, report *domain.AiReport) dto.ReportResponse {
	var content dto.AIReportContent
	if err := json.Unmarshal([]byte(report.Report), &content); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return dto.ReportResponse{}
	}

	return dto.ReportResponse{
		Report: content,
	}
}

func PutReportAsyncHandler(ctx context.HandlerContext,
	aiReportAsyncService service.IAiReportAsyncService,
) dto.CreateReportResponse {
	var request dto.CreateReportRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return dto.CreateReportResponse{}
	}

	if request.Code <= 0 {
		ctx.Status(http.StatusBadRequest)
		return dto.CreateReportResponse{}
	}

	hash, err := aiReportAsyncService.CreateReportRequest(request.Code, request_type.AIREPORT)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return dto.CreateReportResponse{}
	}

	ctx.Status(204)
	return dto.CreateReportResponse{
		Hash:    hash,
		Message: "Запрос принят в работу",
	}
}
