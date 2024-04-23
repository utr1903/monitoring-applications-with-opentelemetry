package server

import (
	"context"
	"net"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/genproto"
	"google.golang.org/grpc"
)

type IServer interface {
	StoreTask(*pb.StoreTaskRequest) *pb.StoreTaskResponse
}

type Server struct {
	logger     logger.ILogger
	grpcServer *grpc.Server
}

type server struct {
	pb.UnimplementedTaskServiceServer

	logger       logger.ILogger
	storeService IStoreService
}

func NewServer(logger logger.ILogger) *Server {
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{
		logger:       logger,
		storeService: NewStoreService(),
	})
	return &Server{
		logger:     logger,
		grpcServer: s,
	}
}

func (s *Server) Run() {
	ctx := context.Background()
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		s.logger.Log(ctx, logger.Error, "Failed to listen.",
			map[string]string{
				"error.message": err.Error(),
			},
		)
	}

	s.logger.Log(ctx, logger.Error, "Server listening...", map[string]string{})
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Log(context.Background(), logger.Error, "Failed to serve.",
			map[string]string{
				"error.message": err.Error(),
			},
		)
	}
}

func (s *server) StoreTask(ctx context.Context, request *pb.StoreTaskRequest) (*pb.StoreTaskResponse, error) {

	result := s.storeService.Store(&StoreRequest{
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

		s.logger.Log(context.Background(), logger.Error, "Storing task failed.",
			map[string]string{
				"error.message": "",
			},
		)
	}

	response := &pb.StoreTaskResponse{
		Code:    code,
		Message: message,
		Task: &pb.Task{
			Id:      result.Task.Id.String(),
			Message: request.Message,
		},
	}

	return response, nil
}
