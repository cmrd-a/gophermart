package domain

import "time"

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
	Accrual    int64
	UploadedAt time.Time
	UserID     int64
}

func NewOrder(number string, status string, accrual int64, uploadedAt time.Time, userID int64) *Order {
	return &Order{
		Number:     number,
		Status:     status,
		Accrual:    accrual,
		UploadedAt: uploadedAt,
		UserID:     userID,
	}
}
