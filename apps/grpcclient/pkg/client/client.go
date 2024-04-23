package client

import (
	"context"

	logger "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"

	pb "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IClient interface {
	StoreTask(ctx context.Context) error
}

type Client struct {
	logger logger.ILogger
	conn   *grpc.ClientConn
	client pb.TaskServiceClient
}

func NewClient(logger logger.ILogger) *Client {
	client := &Client{
		logger: logger,
	}
	return client
}

func (c *Client) Connect(ctx context.Context) error {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Log(ctx, logger.Error, "Connecting to gRPC server failed.",
			map[string]string{
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
		c.logger.Log(ctx, logger.Error, "Storing task failed.",
			map[string]string{
				"error.message": err.Error(),
			})
		return err
	}

	c.logger.Log(ctx, logger.Error, "Storing task suceeded.",
		map[string]string{
			"task.id":      res.GetTask().Id,
			"task.message": res.GetTask().Message,
		})

	return nil
}
