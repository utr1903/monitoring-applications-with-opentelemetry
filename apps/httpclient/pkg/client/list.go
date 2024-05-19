package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	commonhttp "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/http"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
)

func (c *Client) ListTasks(ctx context.Context) error {
	// Create HTTP request
	req, err := c.createListHttpRequest(ctx)
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
	var resBody commonhttp.ListTasksResponse
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
		msg := "Listing tasks failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": resBody.Message,
		})
		return errors.New("call failed")
	}

	c.logger.Log(ctx, loggers.Info, resBody.Message,
		map[string]interface{}{
			"task.count": len(resBody.Body),
		})

	// Add artificial postprocessing step
	c.postprocess(ctx, c.listDelay)
	return nil
}

func (c *Client) createListHttpRequest(ctx context.Context) (*http.Request, error) {
	// Add query params
	u, _ := url.Parse("http://" + c.serverAddress + "/api")
	queryParams := url.Values{}
	queryParams.Set("limit", "5")
	u.RawQuery = queryParams.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		msg := "Creating HTTP request failed."
		c.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": err.Error(),
		})
		return nil, err
	}
	return req, nil
}
