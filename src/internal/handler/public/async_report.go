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

func buildAsyncReportResponse(ctx context.HandlerContext, report *domain.AsyncRequest) dto.ReportResponse {
	var content dto.AIReportContent
	if err := json.Unmarshal([]byte(report.Result), &content); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return dto.ReportResponse{}
	}

	return dto.ReportResponse{
		Report: content,
	}
}

func PostReportAsyncHandler(ctx context.HandlerContext,
	aiReportAsyncService service.IAiReportAsyncService,
) dto.CreateReportResponse {
	var request dto.CreateReportRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return dto.CreateReportResponse{
			Hash:    "",
			Message: err.Error(),
		}
	}

	if request.Code <= 0 {
		ctx.Status(http.StatusBadRequest)
		return dto.CreateReportResponse{
			Hash:    "",
			Message: "code must be positive",
		}
	}

	hash, err := aiReportAsyncService.CreateReportRequest(request.Code, request_type.AIREPORT)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return dto.CreateReportResponse{
			Hash:    "",
			Message: err.Error(),
		}
	}

	ctx.Status(http.StatusAccepted) // дада 202 а не 204 тк 204 это error
	return dto.CreateReportResponse{
		Hash:    hash,
		Message: "Запрос принят в работу",
	}
}

func GetRequestStatusHandler(ctx context.HandlerContext,
	aiReportAsyncService service.IAiReportAsyncService,
) dto.RequestStatusResponse {
	codeStr := ctx.Get("code")
	hash := ctx.Get("hash")

	if codeStr == "" || hash == "" {
		ctx.Status(http.StatusBadRequest)
		return dto.RequestStatusResponse{
			Hash:   "",
			Status: 0,
			Error:  "code and hash parameters are required",
		}
	}

	request, err := aiReportAsyncService.GetRequestStatusByHash(hash)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return dto.RequestStatusResponse{
			Hash:   hash,
			Status: 0,
			Error:  err.Error(),
		}
	}

	if request == nil {
		ctx.Status(http.StatusBadRequest)
		return dto.RequestStatusResponse{
			Hash:   hash,
			Status: 0,
			Error:  "there are no requests with this hash",
		}
	}

	ctx.Status(http.StatusOK)
	return dto.RequestStatusResponse{
		Hash:   request.Hash,
		Status: request.Status,
		Error:  request.Error,
	}

}
