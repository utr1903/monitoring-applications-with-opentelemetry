package server

import (
	"encoding/json"
	"net/http"

	commonhttp "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/http"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
)

func (s *server) writeDeleteResponse(result *services.DeleteResult, w http.ResponseWriter, hasFailed bool) {
	// Create response
	var resBody commonhttp.DeleteTasksResponse
	if hasFailed {
		w.WriteHeader(http.StatusInternalServerError)
		resBody.Message = "Deleting tasks failed."
	} else {
		w.WriteHeader(http.StatusOK)
		resBody.Message = result.Message
	}

	resBodyBytes, _ := json.Marshal(resBody)
	w.Write(resBodyBytes)
}
