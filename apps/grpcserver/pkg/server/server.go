package server

import (
	"context"
	"net"
	"time"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/genproto"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/pkg/config"
	"google.golang.org/grpc"
)

type IServer interface {
	StoreTask(*pb.StoreTaskRequest) *pb.StoreTaskResponse
}

type Server struct {
	port       string
	logger     logger.ILogger
	grpcServer *grpc.Server
}

type server struct {
	pb.UnimplementedTaskServiceServer

	logger logger.ILogger

	storeService services.IStoreService
	storeDelay   int

	listService services.IListService
	listDelay   int

	deleteService services.IDeleteService
	deleteDelay   int
}

func NewServer(cfg *config.Config, logger logger.ILogger) *Server {
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterTaskServiceServer(s, &server{
		logger: logger,

		storeService: services.NewStoreService(),
		storeDelay:   cfg.StoreDelay,

		listService: services.NewListService(),
		listDelay:   cfg.ListDelay,

		deleteService: services.NewDeleteService(),
		deleteDelay:   cfg.DeleteDelay,
	})
	return &Server{
		port:       cfg.Port,
		logger:     logger,
		grpcServer: s,
	}
}

func (s *Server) Run() {
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.logger.Log(ctx, logger.Error, "Failed to listen.",
			map[string]interface{}{
				"error.message": err.Error(),
			},
		)
	}

	s.logger.Log(ctx, logger.Error, "Server listening...", map[string]interface{}{})
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Log(context.Background(), logger.Error, "Failed to serve.",
			map[string]interface{}{
				"error.message": err.Error(),
			},
		)
	}
}

func (s *server) StoreTask(ctx context.Context, request *pb.StoreTaskRequest) (*pb.StoreTaskResponse, error) {
	// Initial artifical delay
	time.Sleep(time.Duration(s.storeDelay) * time.Millisecond)

	// Store task
	result := s.storeService.Store(&services.StoreRequest{
		Task: request.Message,
	})

	var code int32
	var message string

	if result.Result {
		code = 1
		message = "Storing task succeeded."
	} else {
		code = 2
		message = "Storing task failed."

		s.logger.Log(ctx, logger.Error, message, map[string]interface{}{})
	}

	response := &pb.StoreTaskResponse{
		Code:    code,
		Message: message,
		Body: &pb.Task{
			Id:      result.Body.Id.String(),
			Message: request.Message,
		},
	}

	return response, nil
}

func (s *server) ListTasks(ctx context.Context, request *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	// Initial artifical delay
	time.Sleep(time.Duration(s.listDelay) * time.Millisecond)

	// List tasks
	result := s.listService.List(&services.ListRequest{})

	var code int32
	var message string
	if result.Result {
		code = 1
		message = "Listing tasks succeeded."
	} else {
		code = 2
		message = "Listing tasks failed."

		s.logger.Log(ctx, logger.Error, message, map[string]interface{}{})
	}

	tasks := []*pb.Task{}
	for _, task := range result.Body {
		tasks = append(tasks, &pb.Task{
			Id:      task.Id.String(),
			Message: task.Message,
		})
	}

	response := &pb.ListTasksResponse{
		Code:    code,
		Message: message,
		Body:    tasks,
	}

	return response, nil
}

func (s *server) DeleteTasks(ctx context.Context, request *pb.DeleteTasksRequest) (*pb.DeleteTasksResponse, error) {
	// Initial artifical delay
	time.Sleep(time.Duration(s.deleteDelay) * time.Millisecond)

	// Delete tasks
	result := s.deleteService.Delete(&services.DeleteRequest{})

	// Prep response
	var code int32
	if result.Result {
		code = 1
	} else {
		code = 2
		s.logger.Log(ctx, logger.Error, result.Message, map[string]interface{}{})
	}

	response := &pb.DeleteTasksResponse{
		Code:    code,
		Message: result.Message,
	}

	return response, nil
}
