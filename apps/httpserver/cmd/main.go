package main

import (
	"context"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers/logrus"
	otel "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/opentelemetry"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/httpserver/pkg/config"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/httpserver/pkg/server"
)

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
