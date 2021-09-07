package observer

import (
	"github.com/icecolbeveridge/simple_orderbook_2/orderbook"
  	"fmt"	
	"time"
	"encoding/json"
//	"log"
)

const request_time = 5 * time.Second 

// **** OBSERVER

type Observer struct {
	ask_channel chan string
	listening_channel chan string
	candleSticks []orderbook.Candle
	latest_candle_end time.Time
	market_name string
	Server_channel chan string
}

func NewObserver(ob orderbook.Orderbook) Observer {
	var o Observer
	o.ask_channel = ob.Candle_in
	o.listening_channel = ob.Candle_out
	o.Server_channel = make(chan string)
	o.latest_candle_end = time.Now()
	o.market_name = ob.Name 

	go o.sender()
	go o.listener()
	return o 
}

// request candles and send them out

func (o *Observer) sender() {
	for {
		if time.Now().After( o.latest_candle_end.Add(request_time)) {
			cr := orderbook.NewCandleRequest( o.latest_candle_end , request_time)
			 

			o.ask_channel <- cr.ToString()
			time.Sleep(request_time)
		}
	}
}

func (o *Observer) addCandle(c orderbook.Candle) []orderbook.Candle {

	o.candleSticks  = append(o.candleSticks, c)
	return o.candleSticks

}

func (o *Observer) listener() {
	for {
		select  {
		case message := <- o.listening_channel:
			
			var c orderbook.Candle
			_ 				= json.Unmarshal([]byte(message), &c)
			st 				:= orderbook.StringToTime(c.Start_time)
			duration 		:= orderbook.SecondsToDuration(c.Duration_secs)
			end_time 		:= st.Add(duration)
			
			o.candleSticks  = o.addCandle(c)
			
			o.latest_candle_end = end_time
			fmt.Println("sending: " , c.ToString())
			o.Server_channel <- c.ToString()

		}

	}

}

func (o *Observer) RequestCandlesticks() []orderbook.Candle {
	return o.candleSticks
}

func main () {}
