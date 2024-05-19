package server

import (
	"encoding/json"
	"net/http"

	commonhttp "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/http"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
)

func (s *server) writeDeleteResponse(result *services.DeleteResult, w http.ResponseWriter) {
	// Create response
	resBody := &commonhttp.DeleteTasksResponse{
		Message: result.Message,
	}
	resBodyBytes, _ := json.Marshal(resBody)
	w.WriteHeader(http.StatusOK)
	w.Write(resBodyBytes)
}
