package orderbook

import (
	"log"
	"math"
	"time"
)

// **** ORDERBOOK

type Orderbook struct {
	Name string

	trading_channel chan Order
	Candle_in chan string
	Candle_out chan string 

	asks  []Order
	bids  []Order
	execs []Exec

}

// TODO: tidy this up a bit -- I got lost in pointers trying to do polymorphism, but I think I see the way forward now.

type Order interface {
	GetPrice() float64
	GetAmount() float64
	GetSort() string
}

// utility functions

func min(orders []Order) (int, Order) {
	var lowest_order Order
	var lo_index int
	for i, o := range orders {
		if lowest_order == nil {
			lowest_order = o
			lo_index = 0
		} else if o.GetPrice() < lowest_order.GetPrice() && o.GetAmount() > 0. {
			lowest_order = o
			lo_index = i
		}
	}
	return lo_index, lowest_order
}

func max(orders []Order) (int, Order) {
	var out Order
	var hi_index int
	for i , o := range orders {
		if out == nil && o.GetAmount() > 0. {
			out = o
			hi_index = 0
		} else if o.GetPrice() > out.GetPrice() && o.GetAmount() > 0. {
			out = o
			hi_index = i
		}
	}
	return hi_index, out
}

func removeElement(orders []Order, index int) []Order {
	l := len(orders)
	orders[index] = orders[l-1]
	return orders[:l-1]
}

// housekeeping clears the market -- if there's any crossover between bids and asks, match them up and make them Execs.

func (o *Orderbook) housekeep() { // clear the market
	li, lowest_ask := min(o.asks)
	hi, highest_bid := max(o.bids)

	if lowest_ask == nil || highest_bid == nil { return }

	if lowest_ask.GetPrice() < highest_bid.GetPrice() {
		price := 0.5 * (lowest_ask.GetPrice() + highest_bid.GetPrice())
		amount := math.Min(lowest_ask.GetAmount(), highest_bid.GetAmount())
		
		exec := Exec{price: price, amount: amount, timestamp: time.Now()}
		o.execs = append(o.execs, exec)

		new_ask := Ask{price: lowest_ask.GetPrice(), amount: lowest_ask.GetAmount()-amount} 
		new_bid := Bid{price: highest_bid.GetPrice(), amount: highest_bid.GetAmount()-amount}

		if new_ask.GetAmount() > 0 {
			o.asks[li] = &new_ask
		} else {
			o.asks = removeElement(o.asks, li)
		}

		if new_bid.GetAmount() > 0 {
			o.bids[hi] = &new_bid
		} else {
			o.bids = removeElement(o.bids, hi)
		}

	}
}

func (o *Orderbook) generateCandle(request string) Candle {
	return NewCandle( o.execs, request)
}


// listen for orders, put them in the right place

func (o *Orderbook) trade_listener() {
	
	for {
		select {
		case message := <- o.trading_channel:
			switch message.GetSort() {
			case "Ask":
				o.asks = append(o.asks, message)
			case "Bid":
				o.bids = append(o.bids, message)
			default:
				log.Panic("message arrived that wasn't a Bid or Ask")
			}	
		default:
			o.housekeep()
		}

	
	}
}

// listen for candlestick requests and serve them

func (o *Orderbook) candle_listener() {
	for {
		select {
		case request := <- o.Candle_in:
			candle := o.generateCandle(request)
			o.Candle_out <- candle.ToString() 	
		}
	}
}

func NewOrderbook(name string) Orderbook {
	var n Orderbook
	n.Name = name
	n.asks = make( []Order, 0)
	n.bids = make( []Order, 0)
	n.execs = make([]Exec,0)


	n.trading_channel = make(chan Order)
	n.Candle_in = make(chan string)
	n.Candle_out = make(chan string)
	// create channels
	// start a goroutine
	go n.trade_listener()
	go n.candle_listener()
	return n

}

func (o *Orderbook) GetChannel() chan Order {
	return o.trading_channel
}



