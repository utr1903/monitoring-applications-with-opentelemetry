package main

import (
	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/pkg/server"
)

//go:generate go get google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate mkdir -p ../genproto
//go:generate protoc --go_out=../genproto/ --go_opt=paths=source_relative --go-grpc_out=../genproto --go-grpc_opt=paths=source_relative --proto_path=../../proto ../../proto/task.proto

func main() {
	l := logger.NewLogrusLogger("grpcserver")
	s := server.NewServer(l)
	s.Run()
}
