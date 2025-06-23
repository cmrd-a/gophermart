package api

import "time"

type orderStatus string

const (
	REGISTERED orderStatus = "REGISTERED"
	INVALID    orderStatus = "INVALID"
	PROCESSING orderStatus = "PROCESSING"
	PROCESSED  orderStatus = "PROCESSED"
)

type Order struct {
	Number int64 `json:"number" example:"42"`
	// Статус расчёта начисления:
	// * REGISTERED - заказ зарегистрирован, но вознаграждение не рассчитано;
	// * INVALID - заказ не принят к расчёту, и вознаграждение не будет начислено;
	// * PROCESSING - расчёт начисления в процессе;
	// * PROCESSED - расчёт начисления окончен;
	Status     string    `json:"status" example:"PROCESSING" enums:"REGISTERED,INVALID,PROCESSING,PROCESSED"`
	Accural    int64     `json:"accural" example:"500"`
	UploadedAt time.Time `json:"uploaded_at" example:"2025-06-23T23:48:45+03:00"`
}

type Orders []Order
