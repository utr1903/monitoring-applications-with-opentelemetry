package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	commonhttp "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/http"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
)

func (c *Client) DeleteTasks(ctx context.Context) error {
	// Create HTTP request
	req, err := c.createDeleteHttpRequest(ctx)
	if err != nil {
		return err
	}

	// Perform HTTP request
	res, err := c.performHttpRequest(ctx, req)
	if err != nil {
		return err
	}

	// Read HTTP response
	resBodyBytes, err := c.readHttpResponse(ctx, res)
	if err != nil {
		return err
	}

	// Parse HTTP response
	var resBody commonhttp.DeleteTasksResponse
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		msg := "Parsing response body failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})
		return err
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		msg := "Deleting tasks failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": resBody.Message,
		})
		return errors.New("call failed")
	}

	c.logger.Log(ctx, loggers.Info, resBody.Message,
		map[string]interface{}{})

	// Add artificial postprocessing step
	c.postprocess(ctx, c.deleteDelay)
	return nil
}

func (c *Client) createDeleteHttpRequest(ctx context.Context) (*http.Request, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		"http://"+c.serverAddress+"/api",
		nil,
	)
	if err != nil {
		msg := "Writing request body failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})
		return nil, err
	}
	return req, nil
}
