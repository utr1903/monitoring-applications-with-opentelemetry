package server

import (
	"context"
	"log"
	"net"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/genproto"
	"google.golang.org/grpc"
)

type IServer interface {
	StoreTask(*pb.StoreTaskRequest) *pb.StoreTaskResponse
}

type Server struct {
	grpcServer *grpc.Server
}

type server struct {
	pb.UnimplementedTaskServiceServer

	storeService IStoreService
}

func NewServer() *Server {
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{
		storeService: NewStoreService(),
	})
	return &Server{
		grpcServer: s,
	}
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
