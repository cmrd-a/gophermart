package api

import (
	"fmt"
	"time"
)

type UserRegisterRequest struct {
	Login    string `json:"login" example:"user@example.com" binding:"required"`
	Password string `json:"password" example:"password" binding:"required"`
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

type Order struct {
	Number string `json:"number" example:"42"`
	// Статус расчёта начисления:
	// * NEW - заказ загружен в систему, но не попал в обработку;
	// * PROCESSING - вознаграждение за заказ рассчитывается;
	// * INVALID - система расчёта вознаграждений отказала в расчёте;
	// * PROCESSED -  данные по заказу проверены и информация о расчёте успешно получена;
	Status     string   `json:"status" example:"PROCESSING" enums:"NEW,PROCESSING,PROCESSED,INVALID"`
	Accrual    int64    `json:"accrual,omitempty" example:"500"`
	UploadedAt JSONTime `json:"uploaded_at" example:"2025-06-23T23:48:45+03:00"`
}

type Orders []Order
