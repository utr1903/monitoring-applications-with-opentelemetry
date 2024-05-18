package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	commonhttp "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/http"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
)

func (s *server) readListRequestBody(ctx context.Context, bodyReader io.ReadCloser, w http.ResponseWriter) (*commonhttp.ListTasksRequest, error) {
	// Read the request reqBodyBytes
	reqBodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Reading request body failed."
		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})

		// Create response
		resBody := &commonhttp.ListTasksResponse{
			Message: msg,
			Body:    nil,
		}
		resBodyBytes, _ := json.Marshal(resBody)
		w.Write(resBodyBytes)

		return nil, err
	}
	defer bodyReader.Close()

	// Parse the JSON request body
	var reqBody commonhttp.ListTasksRequest
	err = json.Unmarshal(reqBodyBytes, &reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Parsing request body failed."
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
	return &reqBody, nil
}

func (s *server) writeListResponse(result *services.ListResult, w http.ResponseWriter) {
	// Create response
	tasks := []commonhttp.Task{}
	for _, task := range result.Body {
		tasks = append(tasks, commonhttp.Task{
			Id:      task.Id.String(),
			Message: task.Message,
		})
	}

	resBody := &commonhttp.ListTasksResponse{
		Message: result.Message,
		Body:    tasks,
	}
	resBodyBytes, _ := json.Marshal(resBody)
	w.WriteHeader(http.StatusOK)
	w.Write(resBodyBytes)
}
