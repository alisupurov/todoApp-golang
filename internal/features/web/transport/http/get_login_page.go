package web_transport_http

import (
	"net/http"

	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetLoginPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	htmlFile, err := h.webService.GetLoginPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get login.html")
		return
	}

	responseHandler.HTMLResponse(htmlFile)
}
