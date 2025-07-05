package accrual

import (
	"fmt"

	"github.com/cmrd-a/gophermart/internal/config"

	"resty.dev/v3"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetOrderInfo(orderNumber string) (acc OrderInfoResponse, statusCode int, err error) {
	client := resty.New()
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close client: %v\n", closeErr)
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
