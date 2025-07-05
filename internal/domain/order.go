package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	NEW        Status = "NEW"
	PROCESSING Status = "PROCESSING"
	PROCESSED  Status = "PROCESSED"
	INVALID    Status = "INVALID"
)

type Order struct {
	Number     string
	Status     string
	Accrual    decimal.Decimal
	UploadedAt time.Time
	UserID     int64
}

type Balance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}
