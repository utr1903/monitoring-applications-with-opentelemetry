package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/pkg/client"
)

//go:generate go get google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate mkdir -p ../genproto
//go:generate protoc --go_out=../genproto/ --go_opt=paths=source_relative --go-grpc_out=../genproto --go-grpc_opt=paths=source_relative --proto_path=../../proto ../../proto/task.proto

func main() {

	l := logger.NewLogrusLogger("grpcclient")
	c := client.NewClient(l)

	// Connect to grpcserver
	ctx := context.Background()
	err := c.Connect(ctx)
	if err != nil {
		return
	}
	defer c.Close()

	// Wait for signal to shutdown the simulator
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Simulate
	go func() {
		for {
			err = c.StoreTask(ctx)
			if err != nil {
				continue
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			err = c.ListTasks(ctx)
			if err != nil {
				continue
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			err = c.DeleteTasks(ctx)
			if err != nil {
				continue
			}
			time.Sleep(10 * time.Second)
		}
	}()

	<-ctx.Done()
}
