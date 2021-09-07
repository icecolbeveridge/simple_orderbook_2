package trader

import (
	"github.com/icecolbeveridge/simple_orderbook_2/orderbook"
	"math/rand"
	"time"
)

// **** TRADER

type Trader struct {
	name string
	trading_channels []chan orderbook.Order 
}

func (t *Trader) GetName() string {
	return t.name
}

func (t *Trader) GetChannels() []chan orderbook.Order {
	return t.trading_channels
}

func (t *Trader) send_trades() {
	for {
		l := len(t.trading_channels)
		//fmt.Println(l)
		if l > 0 {
			r := rand.Intn(l)
			c := t.trading_channels[r]
			
			if rand.Intn(2) == 0 {
				offer := orderbook.RandomAsk()				
				c <- &offer
			} else {
				offer := orderbook.RandomBid()
				c <- &offer
			}
			
			
		}
		n := time.Duration( rand.Intn(100) + 100)
		time.Sleep(n * time.Millisecond)
		
	}
}

func NewTrader(name string) Trader {
	var t Trader
	t.name = name
	t.trading_channels = make( []chan orderbook.Order, 0)
	return t 
}

func (t *Trader) Run() {
	t.send_trades()
}

func (t *Trader) addChannel(c chan orderbook.Order) []chan orderbook.Order{
	t.trading_channels = append(t.trading_channels, c)
	return t.trading_channels
}

func (t *Trader) Subscribe(o orderbook.Orderbook) Trader {
	t.addChannel(o.GetChannel())
	return *t
}

