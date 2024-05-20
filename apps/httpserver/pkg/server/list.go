package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	commonhttp "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/http"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
)

func (s *server) readListRequestQueryParam(ctx context.Context, r *http.Request, w http.ResponseWriter) (*int64, error) {
	// Read the query limit from request param
	queryParams := r.URL.Query()
	queryLimitParam := queryParams.Get("limit")

	if queryLimitParam == "" {
		queryLimit := int64(5)
		return &queryLimit, nil
	}

	// Parse the query limit param
	queryLimit, err := strconv.ParseInt(queryLimitParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Parsing request query param failed."
		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})

		// Create response
		resBody := &commonhttp.StoreTaskResponse{
			Message: msg,
			Body:    nil,
		}
		resBodyBytes, _ := json.Marshal(resBody)
		w.Write(resBodyBytes)

		return nil, err
	}

	return &queryLimit, nil
}

func (s *server) writeListResponse(result *services.ListResult, w http.ResponseWriter, hasFailed bool) {
	// Create response
	var resBody commonhttp.ListTasksResponse
	if hasFailed {
		w.WriteHeader(http.StatusInternalServerError)
		resBody.Message = "Listing tasks failed."
	} else {
		w.WriteHeader(http.StatusOK)
		tasks := []commonhttp.Task{}
		for _, task := range result.Body {
			tasks = append(tasks, commonhttp.Task{
				Id:      task.Id.String(),
				Message: task.Message,
			})
		}
		resBody.Message = result.Message
		resBody.Body = tasks
	}

	resBodyBytes, _ := json.Marshal(resBody)
	w.Write(resBodyBytes)
}
