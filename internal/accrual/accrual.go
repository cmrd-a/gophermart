package accrual

import (
	"log/slog"

	"github.com/cmrd-a/gophermart/internal/config"

	"resty.dev/v3"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetOrderInfo(orderNumber string) (acc OrderInfoResponse, statusCode int, err error) {
	client := resty.New().SetDebug(true)
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			slog.Warn("Failed to close HTTP client", "error", closeErr)
		}
	}()
	client.SetBaseURL(config.Config.AccrualSystemAddress)

	statusCode = 0
	res, err := client.R().SetPathParam("orderNumber", orderNumber).SetResult(&acc).Get("/api/orders/{orderNumber}")
	if res != nil {
		statusCode = res.StatusCode()
	}
	return acc, statusCode, err
}
