package main

import (
	"context"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/pkg/client"
)

//go:generate go get google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate mkdir -p ../genproto
//go:generate protoc --go_out=../genproto/ --go_opt=paths=source_relative --go-grpc_out=../genproto --go-grpc_opt=paths=source_relative --proto_path=../../proto ../../proto/task.proto

func main() {

	ctx := context.Background()

	l := logger.NewLogrusLogger("grpcclient")
	c := client.NewClient(l)
	err := c.Connect(ctx)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.StoreTask(ctx)
	if err != nil {
		return
	}

	err = c.ListTasks(ctx)
	if err != nil {
		return
	}
}
