package accrual

type Status string

const (
	REGISTERED Status = "REGISTERED" //заказ зарегистрирован, но вознаграждение не рассчитано
	INVALID    Status = "INVALID"    //заказ не принят к расчёту, и вознаграждение не будет начислено
	PROCESSING Status = "PROCESSING" //расчёт начисления в процессе
	PROCESSED  Status = "PROCESSED"  //расчёт начисления окончен
)

type OrderInfoResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
