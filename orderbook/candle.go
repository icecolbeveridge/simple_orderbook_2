package orderbook

import (
	"time"
	"encoding/json"
	"log"
)

const time_format = time.RFC1123


// utility functions

func StringToTime(s string) time.Time {
	t, _ :=  time.Parse(time_format, s)
	return t 
}

func TimeToString(t time.Time) string {
	return t.Format(time_format)
}

func DurationToSeconds(d time.Duration) float64 {
	return d.Seconds()
}

func SecondsToDuration(i float64) time.Duration {
	return time.Duration(i) * time.Second
}


// **** CANDLEREQUEST

type CandleRequest struct {
	Start_time string 		  // represents time.Time but needs to be jsonable
	Duration_secs  float64    // do. time.Duration
}

func NewCandleRequest( start_time time.Time, duration time.Duration) CandleRequest {
	return CandleRequest{ Start_time:TimeToString(start_time), Duration_secs: DurationToSeconds(duration) }
}

// **** CANDLE

type Candle struct {
	Start_time string 			// needs to be a string for json
	Duration_secs float64       // same principle
	Open float64
	Close float64
	High float64
	Low float64
	Volume float64
}


func NewCandle( execs []Exec, request string) Candle {
	var cr CandleRequest
	err := json.Unmarshal([]byte(request), &cr) 
	if err != nil {log.Panic(err)}
	empty := true
	open := 0.
	close := 0.
	high := 0.
	low := 0.
	volume := 0.
	var early time.Time
	var late time.Time
	start_time := StringToTime(cr.Start_time)
	duration   := SecondsToDuration(cr.Duration_secs)

	for _, exec := range execs {
		t := exec.timestamp
		if t.Before( start_time.Add(duration)) && t.After( start_time) {
			if empty {
				empty = false
				early = t
				late = t
				open = exec.price
				close = exec.price
				high = exec.price
				low = exec.price
				volume = exec.amount
			}
			if t.Before(early) {
				early = t
				open = exec.price
			} else if t.After(late) {
				late = t
				close = exec.price
			}
			if exec.price > high {
				high = exec.price
			} else if exec.price < low {
				low = exec.price
			}
			volume += exec.amount
		}
	}
	return Candle{ Start_time: cr.Start_time, Duration_secs: cr.Duration_secs, Open: open, Close: close, High: high, Low: low, Volume: volume}

}

// json marshalling

func (c Candle) ToString() string {
	j, _ := json.Marshal(c)
	return string(j)
}

func (c CandleRequest) ToString() string {
	j, err := json.Marshal(c)
	if err != nil {
		log.Panic(err)
	}
	return string(j)
}

//func main() {}