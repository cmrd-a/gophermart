package domain

import "time"

type Order struct {
	Number     string
	Status     string
	Accural    int64
	UploadedAt time.Time
	UserID     int64
}

func NewOrder(number string, status string, accural int64, uploadedAt time.Time, userID int64) *Order {
	return &Order{
		Number:     number,
		Status:     status,
		Accural:    accural,
		UploadedAt: uploadedAt,
		UserID:     userID,
	}
}
