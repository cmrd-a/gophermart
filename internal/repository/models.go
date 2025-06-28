package repository

import "time"

type Order struct {
	Number     string
	Status     string
	Accural    int64
	UploadedAt time.Time
	UserID     int64
}
