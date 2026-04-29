package public

import (
	"backend/src/internal/async/request"
	context "backend/src/internal/context/abstract"
	service "backend/src/internal/service/abstract"
	"net/http"
	"strconv"
)

func GetReportAsyncHandler(ctx context.HandlerContext, aiReportAsyncService service.IAiReportAsyncService) {
	codeStr := ctx.Get("code")

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	report, err := aiReportAsyncService.GetReportByCode(request.AIREPORT)
	if err != nil || report == nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

}
