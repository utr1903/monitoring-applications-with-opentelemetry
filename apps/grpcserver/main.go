package main

import (
	"context"
	"log"
	"net"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/genproto/task"
	"google.golang.org/grpc"
)

//go:generate go get google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate protoc --go_out=./ --go-grpc_out=./ --proto_path=../proto ../proto/task.proto

type server struct {
	pb.UnimplementedTaskServiceServer
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SayHello implements TaskServer
func (s *server) StoreTask(ctx context.Context, in *pb.StoreTaskRequest) (*pb.StoreTaskResponse, error) {
	log.Printf("Received: %v", in.GetTask())
	return &pb.StoreTaskResponse{
		Result:  "OK",
		Message: "Hello " + in.GetTask()}, nil
}
