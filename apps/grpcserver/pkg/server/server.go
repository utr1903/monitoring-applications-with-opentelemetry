package server

import (
	"context"
	"net"
	"time"

	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	services "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/services"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/genproto"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/pkg/config"
	"google.golang.org/grpc"
)

type Server struct {
	port       string
	logger     loggers.ILogger
	grpcServer *grpc.Server
}

type server struct {
	pb.UnimplementedTaskServiceServer

	logger loggers.ILogger

	storeService services.IStoreService
	storeDelay   int

	listService services.IListService
	listDelay   int

	deleteService services.IDeleteService
	deleteDelay   int
}

func NewServer(cfg *config.Config, logger loggers.ILogger) *Server {
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterTaskServiceServer(s, &server{
		logger: logger,

		storeService: services.NewStoreService(logger, cfg.CreateDbNotReachableError),
		storeDelay:   cfg.StoreDelay,

		listService: services.NewListService(logger, cfg.CreateDbNotReachableError),
		listDelay:   cfg.ListDelay,

		deleteService: services.NewDeleteService(logger, cfg.CreateDbNotReachableError),
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
		s.logger.Log(ctx, loggers.Error, "Failed to listen.",
			map[string]interface{}{
				"error.message": err.Error(),
			},
		)
	}

	s.logger.Log(ctx, loggers.Error, "Server listening...", map[string]interface{}{})
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Log(context.Background(), loggers.Error, "Failed to serve.",
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
	result, err := s.storeService.Store(ctx, &services.StoreRequest{
		Task: request.Message,
	})
	if err != nil {
		return nil, err
	}

	return &pb.StoreTaskResponse{
		Message: result.Message,
		Body: &pb.Task{
			Id:      result.Body.Id.String(),
			Message: result.Body.Message,
		}}, err
}

func (s *server) ListTasks(ctx context.Context, request *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	// Initial artifical delay
	time.Sleep(time.Duration(s.listDelay) * time.Millisecond)

	// List tasks
	result, err := s.listService.List(ctx, &services.ListRequest{})
	if err != nil {
		return nil, err
	}

	tasks := []*pb.Task{}
	for _, task := range result.Body {
		tasks = append(tasks, &pb.Task{
			Id:      task.Id.String(),
			Message: task.Message,
		})
	}

	return &pb.ListTasksResponse{
		Message: result.Message,
		Body:    tasks,
	}, err
}

func (s *server) DeleteTasks(ctx context.Context, request *pb.DeleteTasksRequest) (*pb.DeleteTasksResponse, error) {
	// Initial artifical delay
	time.Sleep(time.Duration(s.deleteDelay) * time.Millisecond)

	// Delete tasks
	result, err := s.deleteService.Delete(ctx, &services.DeleteRequest{})

	if err != nil {
		return nil, err
	}

	return &pb.DeleteTasksResponse{
		Message: result.Message,
	}, nil
}
