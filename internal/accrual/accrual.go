package accrual

import (
	"fmt"

	"resty.dev/v3"
)

func GetAccrual(orderNumber string) (AccrualStatus, int64, error) {
	client := resty.New()
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close client: %v\n", closeErr)
		}
	}()
	client.SetBaseURL("http://localhost:8080")

	res, err := client.R().SetPathParam("orderNumber", orderNumber).Get("/api/orders/{orderNumber}")
	fmt.Println(res.StatusCode())
	return UNSPECIFIED, 0, err
}
