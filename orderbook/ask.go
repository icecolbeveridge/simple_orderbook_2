package orderbook

import (
	"math/rand"
)

// **** ASK

type Ask struct {
	sort string
	price float64
	amount float64
}

func (a *Ask) GetPrice() float64 { return a.price} // should be pointerless
func (a *Ask) GetAmount() float64 { return a.amount}
func (a *Ask) GetSort() string { return a.sort}
func (a *Ask) ReduceAmount(n float64) Ask { return Ask{price: a.price, amount: a.amount - n}} 

func RandomAsk() Ask {
	return Ask{price:100 + 20 * rand.Float64(), amount: 1. + float64(rand.Intn(10)), sort:"Ask"}
}