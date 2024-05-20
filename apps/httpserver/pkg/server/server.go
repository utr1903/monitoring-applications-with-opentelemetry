package server

import (
	"context"
	"net/http"
	"time"

	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/httpserver/pkg/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Server struct {
	logger     loggers.ILogger
	httpServer *http.Server
}

type server struct {
	logger loggers.ILogger

	storeService services.IStoreService
	storeDelay   int

	listService services.IListService
	listDelay   int

	deleteService services.IDeleteService
	deleteDelay   int
}

func NewServer(cfg *config.Config, logger loggers.ILogger) *Server {

	server := &server{
		logger: logger,

		storeService: services.NewStoreService(logger, cfg.CreateDbNotReachableError),
		storeDelay:   cfg.StoreDelay,

		listService: services.NewListService(logger, cfg.CreateDbNotReachableError),
		listDelay:   cfg.ListDelay,

		deleteService: services.NewDeleteService(logger, cfg.CreateDbNotReachableError),
		deleteDelay:   cfg.DeleteDelay,
	}

	mux := http.NewServeMux()
	mux.Handle("/api", otelhttp.NewHandler(http.HandlerFunc(server.Handle), "api"))

	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: mux,
		},
	}
}

func (s *Server) Run() {
	ctx := context.Background()
	s.logger.Log(ctx, loggers.Info, "Server listening...", map[string]interface{}{})
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.logger.Log(ctx, loggers.Error, "Failed to serve.",
			map[string]interface{}{
				"error.message": err.Error(),
			},
		)
	}
}

func (s *server) Handle(w http.ResponseWriter, r *http.Request) {
	s.logger.Log(r.Context(), loggers.Info, "Handling request...", map[string]interface{}{})

	switch {
	case r.Method == http.MethodPost:
		s.StoreTask(w, r)
	case r.Method == http.MethodGet:
		s.ListTasks(w, r)
	case r.Method == http.MethodDelete:
		s.DeleteTasks(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	s.logger.Log(r.Context(), loggers.Info, "Request is handled.", map[string]interface{}{})
}

func (s *server) StoreTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initial artifical delay
	time.Sleep(time.Duration(s.storeDelay) * time.Millisecond)

	// Read the request body
	reqBody, err := s.readStoreRequestBody(ctx, r.Body, w)
	if err != nil {
		return
	}

	// Store task
	result, err := s.storeService.Store(ctx, &services.StoreRequest{
		Task: reqBody.Message,
	})

	// Write response
	s.writeStoreResponse(result, w, err != nil)
}

func (s *server) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initial artifical delay
	time.Sleep(time.Duration(s.listDelay) * time.Millisecond)

	// Read the request body
	queryLimit, err := s.readListRequestQueryParam(ctx, r, w)
	if err != nil {
		return
	}

	// List tasks
	result, err := s.listService.List(ctx, &services.ListRequest{
		Limit: *queryLimit,
	})

	// Write response
	s.writeListResponse(result, w, err != nil)
}

func (s *server) DeleteTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initial artifical delay
	time.Sleep(time.Duration(s.deleteDelay) * time.Millisecond)

	// Delete tasks
	result, err := s.deleteService.Delete(ctx, &services.DeleteRequest{})

	// Write response
	s.writeDeleteResponse(result, w, err != nil)
}
