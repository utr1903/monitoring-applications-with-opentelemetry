package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	otel "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/opentelemetry"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/pkg/client"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/pkg/config"
)

//go:generate go get google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate mkdir -p ../genproto
//go:generate protoc --go_out=../genproto/ --go_opt=paths=source_relative --go-grpc_out=../genproto --go-grpc_opt=paths=source_relative --proto_path=../../proto ../../proto/task.proto

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogrusLogger(cfg.ServiceName)

	// Get context
	ctx := context.Background()

	// Create tracer provider
	tp := otel.NewTraceProvider(ctx)
	defer otel.ShutdownTraceProvider(ctx, tp)

	// Create metric provider
	mp := otel.NewMetricProvider(ctx)
	defer otel.ShutdownMetricProvider(ctx, mp)

	// Collect runtime metrics
	otel.StartCollectingRuntimeMetrics()

	clt := client.NewClient(log)

	// Connect to grpcserver
	err := clt.Connect(ctx)
	if err != nil {
		return
	}
	defer clt.Close()

	// Wait for signal to shutdown the simulator
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Simulate
	go func() {
		for {
			err = clt.StoreTask(ctx)
			if err != nil {
				continue
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			err = clt.ListTasks(ctx)
			if err != nil {
				continue
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			err = clt.DeleteTasks(ctx)
			if err != nil {
				continue
			}
			time.Sleep(10 * time.Second)
		}
	}()

	<-ctx.Done()
}
