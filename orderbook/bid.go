package orderbook

import (
	"math/rand"
)

// BID

type Bid struct {
	price float64
	amount float64
	sort string
}

func (a *Bid) GetPrice() float64 { return a.price} // these should all be pointerless so I can use the Order interface
func (a *Bid) GetAmount() float64 { return a.amount}
func (a *Bid) GetSort() string { return a.sort}
func (a *Bid) ReduceAmount(n float64) Bid { return Bid{price: a.price, amount: a.amount - n} } 


func RandomBid() Bid {
	return Bid{price:100 + 20 * rand.Float64(), amount: 1. + float64(rand.Intn(10)), sort:"Bid"}
}