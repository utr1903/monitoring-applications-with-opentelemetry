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

func (s *server) readDeleteRequestBody(ctx context.Context, bodyReader io.ReadCloser, w http.ResponseWriter) (*commonhttp.DeleteTasksRequest, error) {
	// Read the request reqBodyBytes
	reqBodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Reading request body failed."
		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})

		// Create response
		resBody := &commonhttp.DeleteTasksResponse{
			Message: msg,
		}
		resBodyBytes, _ := json.Marshal(resBody)
		w.Write(resBodyBytes)

		return nil, err
	}
	defer bodyReader.Close()

	// Parse the JSON request body
	var reqBody commonhttp.DeleteTasksRequest
	err = json.Unmarshal(reqBodyBytes, &reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Parsing request body failed."
		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})

		// Create response
		resBody := &commonhttp.DeleteTasksResponse{
			Message: msg,
		}
		resBodyBytes, _ := json.Marshal(resBody)
		w.Write(resBodyBytes)

		return nil, err
	}
	return &reqBody, nil
}

func (s *server) writeDeleteResponse(result *services.DeleteResult, w http.ResponseWriter) {
	// Create response
	resBody := &commonhttp.DeleteTasksResponse{
		Message: result.Message,
	}
	resBodyBytes, _ := json.Marshal(resBody)
	w.Write(resBodyBytes)
}
