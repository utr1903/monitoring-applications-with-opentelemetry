package client

import (
	"context"
	"errors"
	"time"

	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/genproto"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	logger loggers.ILogger

	serviceName   string
	serverAddress string

	conn   *grpc.ClientConn
	client pb.TaskServiceClient

	storeDelay  int
	listDelay   int
	deleteDelay int

	createPostprocessingError bool
	createPostprocessingDelay bool
}

func NewClient(cfg *config.Config, logger loggers.ILogger) *Client {
	client := &Client{
		logger: logger,

		serverAddress: cfg.ServerAddress,
		serviceName:   cfg.ServiceName,

		storeDelay:  cfg.StoreDelay,
		listDelay:   cfg.ListDelay,
		deleteDelay: cfg.DeleteDelay,

		createPostprocessingError: cfg.CreatePostprocessingError,
		createPostprocessingDelay: cfg.CreatePostprocessingDelay,
	}
	return client
}

func (c *Client) Connect(ctx context.Context) error {
	conn, err := grpc.Dial(c.serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		c.logger.Log(ctx, loggers.Error, "Connecting to gRPC server failed.",
			map[string]interface{}{
				"error.message": err.Error(),
			})
		return err
	}
	c.conn = conn
	c.client = pb.NewTaskServiceClient(conn)

	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) StoreTask(ctx context.Context) error {
	res, err := c.client.StoreTask(ctx, &pb.StoreTaskRequest{
		Message: "Some task.",
	})
	if err != nil {
		c.logger.Log(ctx, loggers.Error, "Storing task failed.",
			map[string]interface{}{
				"error.message": err.Error(),
			})
		return err
	}

	c.logger.Log(ctx, loggers.Info, "Storing task suceeded.",
		map[string]interface{}{
			"task.id":      res.GetBody().Id,
			"task.message": res.GetBody().Message,
		})

	// Add artificial postprocessing step
	c.postprocess(ctx, c.storeDelay)
	return nil
}

func (c *Client) ListTasks(ctx context.Context) error {
	res, err := c.client.ListTasks(ctx, &pb.ListTasksRequest{})
	if err != nil {
		c.logger.Log(ctx, loggers.Error, "Listing task failed.",
			map[string]interface{}{
				"error.message": err.Error(),
			})
		return err
	}

	c.logger.Log(ctx, loggers.Info, "Listing task suceeded.",
		map[string]interface{}{
			"task.count": len(res.GetBody()),
		})

	// Add artificial postprocessing step
	c.postprocess(ctx, c.listDelay)
	return nil
}

func (c *Client) DeleteTasks(ctx context.Context) error {
	_, err := c.client.DeleteTasks(ctx, &pb.DeleteTasksRequest{})
	if err != nil {
		c.logger.Log(ctx, loggers.Error, "Deleting task failed.",
			map[string]interface{}{
				"error.message": err.Error(),
			})
		return err
	}

	c.logger.Log(ctx, loggers.Info, "Deleting task suceeded.",
		map[string]interface{}{})

	// Add artificial postprocessing step
	c.postprocess(ctx, c.deleteDelay)
	return nil
}

func (c *Client) postprocess(ctx context.Context, duration int) {
	// Get current span
	parentSpan := trace.SpanFromContext(ctx)

	c.logger.Log(ctx, loggers.Info, "Postprocessing...",
		map[string]interface{}{})

	// Create postprocessing span
	_, span := parentSpan.TracerProvider().
		Tracer(c.serviceName).
		Start(
			ctx,
			"postprocessing",
			trace.WithSpanKind(trace.SpanKindInternal),
		)
	defer span.End()

	if c.createPostprocessingError {
		err := errors.New("could not find postprocessing schema")
		span.SetStatus(codes.Error, "Postprocessing failed.")
		span.RecordError(err)

		c.logger.Log(ctx, loggers.Error, "Postprocessing failed.",
			map[string]interface{}{
				"error.message": "Postprocessing step crashed due to singularity in calculation.",
			})

		return
	}

	if c.createPostprocessingDelay {
		c.logger.Log(ctx, loggers.Warning, "Postprocessing will take longer.",
			map[string]interface{}{
				"error.message": "Postprocessing schema cache could not be found. Calculating from scratch.",
			})
		time.Sleep(time.Second)
		return
	}

	time.Sleep(time.Duration(duration) * time.Microsecond)
}
