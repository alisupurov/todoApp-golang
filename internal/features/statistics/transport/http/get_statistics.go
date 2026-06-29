package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated             int      `json:"tasks_created"               example:"10"`
	TasksCompleted           int      `json:"tasks_completed"             example:"6"`
	TaskCompletedRate        *float64 `json:"tasks_completed_rate"        example:"60"`
	TasksAverageCompleteTime *string  `json:"tasks_average_complete_time" example:"9h59m18s"`
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompleteTime != nil {
		duration := statistics.TasksAverageCompleteTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated: statistics.TasksCreated,
		TasksCompleted: statistics.TasksCompleted,
		TaskCompletedRate: statistics.TaskCompletedRate,
		TasksAverageCompleteTime: avgTime,
	}
}

// GetStatistics godoc
// @Summary      Получить статистику по задачам
// @Description  Получить агрегированную статистику по задачам с фильтрами по автору и периоду
// @Tags         statistics
// @Produce      json
// @Param        user_id query int false "ID автора задачи"
// @Param        from query string false "Дата начала периода (YYYY-MM-DD)"
// @Param        to query string false "Дата конца периода (YYYY-MM-DD)"
// @Success      200  {object}  GetStatisticsResponse "Статистика"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	params, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return 
	}

	domain, err := h.statisticsService.GetStatistics(ctx, params.userId, params.from, params.to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)

		return
	}

	response := toDTOFromDomain(domain)

	responseHandler.JSONResponse(response, http.StatusOK)
}

type QueryParams struct {
	userId *int
	from *time.Time
	to *time.Time
}

func getUserIDFromToQueryParams(r *http.Request) (QueryParams, error) {
	const (
		userdIdQueryParam = "user_id"
		FromQueryParam = "from"
		ToQueryParam = "to"
	)

	userId, err := core_http_request.GetIntQueryParam(r, userdIdQueryParam)
	if err != nil {
		return QueryParams{}, fmt.Errorf("failed to get user_id query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, FromQueryParam)
	if err != nil {
		return QueryParams{}, fmt.Errorf("failed to get from query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, ToQueryParam)
	if err != nil {
		return QueryParams{}, fmt.Errorf("failed to get to query param: %w", err)
	}

	return QueryParams{
		userId: userId,
		from: from,
		to: to,
	}, nil
}