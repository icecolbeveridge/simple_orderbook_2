package orderbook

import (
	"time"
)

// **** EXEC

type Exec struct {
	price float64
	amount float64
	timestamp time.Time
}

func (a *Exec) GetPrice() float64 { return a.price}
func (a *Exec) GetAmount() float64 { return a.amount}

