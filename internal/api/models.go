package api

import (
	"fmt"
	"time"
)

type UserRegisterRequest struct {
	Login    string `json:"login" example:"user@example.com" binding:"required"`
	Password string `json:"password" example:"password" binding:"required"`
}

type orderStatus string

const (
	REGISTERED orderStatus = "REGISTERED"
	INVALID    orderStatus = "INVALID"
	PROCESSING orderStatus = "PROCESSING"
	PROCESSED  orderStatus = "PROCESSED"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

type Order struct {
	Number string `json:"number" example:"42"`
	// Статус расчёта начисления:
	// * REGISTERED - заказ зарегистрирован, но вознаграждение не рассчитано;
	// * INVALID - заказ не принят к расчёту, и вознаграждение не будет начислено;
	// * PROCESSING - расчёт начисления в процессе;
	// * PROCESSED - расчёт начисления окончен;
	Status     string   `json:"status" example:"PROCESSING" enums:"REGISTERED,INVALID,PROCESSING,PROCESSED"`
	Accural    int64    `json:"accural,omitempty" example:"500"`
	UploadedAt JSONTime `json:"uploaded_at" example:"2025-06-23T23:48:45+03:00"`
}

type Orders []Order
