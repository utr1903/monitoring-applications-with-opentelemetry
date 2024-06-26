package main

import (
	"context"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	otel "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/opentelemetry"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/pkg/config"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcserver/pkg/server"
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

	srv := server.NewServer(cfg, log)
	srv.Run()
}
