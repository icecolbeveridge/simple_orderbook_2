/* ---------
/. Simulate multiple orderbooks with concurrent traders
/. Also run a server exposing latest candlesticks from each market
/. Starts a server on port 8642
/. visit http://localhost:8642/BIGCO, http://localhost:8642/MIDCO or http://localhost:8642/WEECO to see the five-second candlesticks for each market.
/* --------- */


package main

import (
	"github.com/icecolbeveridge/simple_orderbook_2/orderbook" 	
	"github.com/icecolbeveridge/simple_orderbook_2/trader" 		
	
	"github.com/icecolbeveridge/simple_orderbook_2/observer"    

	"net/http"
	"fmt"
	"log"
	"bufio"
	"os"
)


func writeToFile(message string, filename string) {
	// append message to file with given name
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if _, err = f.WriteString(message + "\n"); err != nil {
        panic(err)
    }
}

func candle_listener(n string, c chan string, filename string) {
	// receives candlesticks from observers

	for {
		select {
		case message := <- c:
			
			writeToFile(string(message) + "\n", filename)
		
		}
	}

}


func main() {

	trader_names 	:= []string{ "Andy", "Beth", "Chad"}
	traders 		:= make(map[string] trader.Trader)

	market_names 	:= []string{ "BIGCO", "MIDCO", "WEECO"}
	
	observers 		:= make(map[string] observer.Observer )

	// create traders
	for _, t := range trader_names {
		traders[t] = trader.NewTrader(t) // TODO (should automatically run a goroutine generating orders on each of its orderbook channels)
	}

	
	
	

	for _, n := range market_names {

		// create orderbooks

		o := orderbook.NewOrderbook(n) 

		// subscribe traders to orderbooks

		for _, t := range trader_names {
			trader := traders[t]
			traders[t] = trader.Subscribe(o) 
			
		}

		// create observers (who handle the candles)

		observers[n] = observer.NewObserver(o)

		// create file to hold candlestick data (this is as much to practice IO as anything)

		filename  := "./" + n + ".data"
		_,err := os.Create(filename)
		if err != nil {
			log.Panic(err)
		}

		// listen for candles sent by the observers

		go candle_listener(n, observers[n].Server_channel, filename)
	}

	// let the traders trade

	for _, t := range trader_names {
		trader := traders[t]
		go trader.Run()
	}

	// create server and serve information

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
                        market := r.URL.Path[1:]

                        // read data from the relevant path
                        file, err := os.Open("./" + market + ".data")
                        defer file.Close()
                        
                        if err != nil {
                        	fmt.Fprintf(w, "No file found")
                        	
                        }

                    	// read data from file    
        			                
                	    scanner := bufio.NewScanner(file)
						for scanner.Scan() {
	                        fmt.Fprintf(w, scanner.Text() )      
    	                } 
        	            

                	}) 

    log.Fatal(http.ListenAndServe(":8642", nil))

}
