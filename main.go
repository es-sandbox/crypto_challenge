package main

import (
	"time"
	"math/rand"
	"os"
	"os/signal"
	"fmt"
	"net/http"
	"sync"
	"encoding/json"
	"sort"
)

var (
	exchange *Exchange
	dynamicPriceEstimator priceEstimator

	dealChan chan *deal
	archiveDeal []*deal
	archiveDealMtx *sync.Mutex
)

// NOTE: must be run as goroutine
func enableOrderGeneration(datasetSource *datasetSource, orderChan chan <- *order, timeout time.Duration) {
	// TODO(evg): using time.Ticker instead of sleep
	for {
		var order *order
		if rand.Intn(2) == 0 {
			order = datasetSource.getNextFilteredSellOrder()
		} else {
			order = datasetSource.getNextFilteredBuyOrder()
		}
		orderChan <- order

		time.Sleep(timeout)
	}
}

func priceHadnler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		price := dynamicPriceEstimator.GetPrice()
		data, err := json.Marshal(price)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(data); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func buyOrderHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		exchange.activeOrderMtx.Lock()
		defer exchange.activeOrderMtx.Unlock()

		data, err := json.Marshal(exchange.activeBuyOrderSlice)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(data); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("GET /order/buy: OK")
	}
}

func sellOrderHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		exchange.activeOrderMtx.Lock()
		defer exchange.activeOrderMtx.Unlock()

		less := func(i, j int) bool {
			return exchange.activeSellOrderSlice[i].Price < exchange.activeSellOrderSlice[j].Price
		}
		sort.Slice(exchange.activeSellOrderSlice, less)

		data, err := json.Marshal(exchange.activeSellOrderSlice)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(data); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("GET /order/sell: OK")
	}
}

func dealsHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case "GET":
			// fmt.Println(archiveDealMtx)
			// fmt.Println(archiveDeal)
			// return

			archiveDealMtx.Lock()
			data, err := json.Marshal(archiveDeal)
			if err != nil {
				fmt.Println(err)
				return
			}
			archiveDealMtx.Unlock()

			resp.Header().Set("Access-Control-Allow-Origin", "*")

			if _, err := resp.Write(data); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("GET /deals: OK")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// staticPriceEstimator := newStaticPriceEstimator(1e4)
	dynamicPriceEstimator = newDynamicPriceEstimator(1e4)
	staticAmountGenerator := newStaticAmountGenerator(1e4, 1)
	datasetSource := newDatasetSource(dynamicPriceEstimator, staticAmountGenerator)

	orderChanBufferSize := 20
	orderChan := make(chan *order, orderChanBufferSize)
	go enableOrderGeneration(datasetSource, orderChan, time.Millisecond * 500)

	dealChanBufferSize := 20
	dealChan = make(chan *deal, dealChanBufferSize)
	exchange = newExchange(orderChan, dealChan)
	go exchange.start()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	go func(){
		for sig := range shutdownChan {
			if sig == os.Interrupt {
				fmt.Println("Shutdown")
				os.Exit(0)
			}
			// sig is a ^C, handle it
		}
	}()

	archiveDeal = make([]*deal, 0)
	archiveDealMtx = &sync.Mutex{}
	go func() {
		for {
			deal := <-dealChan

			archiveDealMtx.Lock()
			archiveDeal = append(archiveDeal, deal)
			archiveDealMtx.Unlock()

			fmt.Println(deal)
		}
	}()

	go exchange.enableAnalyze(time.Millisecond * 1000)

	http.HandleFunc("/deals", dealsHandler)
	http.HandleFunc("/order/sell", sellOrderHandler)
	http.HandleFunc("/order/buy", buyOrderHandler)
	http.HandleFunc("/price", priceHadnler)

	listenAddr := "localhost:9000"
	fmt.Printf("listen: %v\n", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}