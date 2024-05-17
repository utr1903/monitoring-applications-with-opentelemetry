package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/opentelemetry"
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
	tp := opentelemetry.NewTraceProvider(ctx)
	defer opentelemetry.ShutdownTraceProvider(ctx, tp)

	// Create metric provider
	mp := opentelemetry.NewMetricProvider(ctx)
	defer opentelemetry.ShutdownMetricProvider(ctx, mp)

	// Collect runtime metrics
	opentelemetry.StartCollectingRuntimeMetrics()

	clt := client.NewClient(cfg, log)

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
			ctx, span := createRootSpan(cfg.ServiceName, ctx)
			clt.StoreTask(ctx)
			span.End()

			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			ctx, span := createRootSpan(cfg.ServiceName, ctx)
			clt.ListTasks(ctx)
			span.End()

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			ctx, span := createRootSpan(cfg.ServiceName, ctx)
			clt.DeleteTasks(ctx)
			span.End()

			time.Sleep(10 * time.Second)
		}
	}()

	<-ctx.Done()
}

func createRootSpan(serviceName string, ctx context.Context) (context.Context, trace.Span) {
	// Create root span
	ctx, span := otel.GetTracerProvider().
		Tracer(serviceName).
		Start(
			ctx,
			"root",
			trace.WithSpanKind(trace.SpanKindServer),
		)

	return ctx, span
}
