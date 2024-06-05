package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/httpclient/pkg/config"
)

type Client struct {
	logger loggers.ILogger

	serviceName   string
	serverAddress string

	client *http.Client

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

		client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},

		storeDelay:  cfg.StoreDelay,
		listDelay:   cfg.ListDelay,
		deleteDelay: cfg.DeleteDelay,

		createPostprocessingError: cfg.CreatePostprocessingError,
		createPostprocessingDelay: cfg.CreatePostprocessingDelay,
	}
	return client
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
		err := errors.New("crashed due to singularity")
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

func (c *Client) performHttpRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Perform HTTP request
	c.logger.Log(ctx, loggers.Info, "Performing the HTTP request...", map[string]interface{}{})
	res, err := c.client.Do(req)
	if err != nil {
		msg := "Performing HTTP request failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})
		return nil, err
	}

	return res, nil
}

func (c *Client) readHttpResponse(ctx context.Context, res *http.Response) ([]byte, error) {
	// Read HTTP response
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		msg := "Reading response body failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})
		return nil, err
	}
	defer res.Body.Close()

	return resBodyBytes, nil
}
