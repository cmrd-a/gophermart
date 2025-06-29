package accrual

type AccrualStatus string

const (
	UNSPECIFIED AccrualStatus = "UNSPECIFIED" //статус не определён
	REGISTERED  AccrualStatus = "REGISTERED"  //заказ зарегистрирован, но вознаграждение не рассчитано
	INVALID     AccrualStatus = "INVALID"     //заказ не принят к расчёту, и вознаграждение не будет начислено
	PROCESSING  AccrualStatus = "PROCESSING"  //расчёт начисления в процессе
	PROCESSED   AccrualStatus = "PROCESSED"   //расчёт начисления окончен
)

type Accrual struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual,omitempty"`
}
