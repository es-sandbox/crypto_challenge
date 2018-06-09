package main

import (
	"log"
	"sync"
	"fmt"
	"sort"
	"time"
	"math"
)

type Exchange struct {
	orderChan <-chan *order
	dealChan  chan *deal

	activeOrderMtx       *sync.Mutex
	activeSellOrderSlice []*order
	activeBuyOrderSlice  []*order
}

func newExchange(orderChan <-chan *order, dealChan chan *deal) *Exchange {
	return &Exchange{
		orderChan: orderChan,
		dealChan: dealChan,

		activeOrderMtx:       &sync.Mutex{},
		activeSellOrderSlice: make([]*order, 0),
		activeBuyOrderSlice:  make([]*order, 0),
	}
}

// NOTE: must be run as goroutine
func (e *Exchange) start() {
	for {
		order := <-e.orderChan

		e.activeOrderMtx.Lock()
		switch order.OrderType {
		case sellOrderType:
			e.activeSellOrderSlice = append(e.activeSellOrderSlice, order)
			fmt.Println("+")
		case buyOrderType:
			e.activeBuyOrderSlice = append(e.activeBuyOrderSlice, order)
			fmt.Println("-")
		default:
			log.Println("[exchange]: <Unknown order type>")
		}
		e.activeOrderMtx.Unlock()

		// log.Println(order)
	}
}

// NOTE: must be run as goroutine
func (e *Exchange) enableAnalyze(timeout time.Duration) {
	for {
		e.analyzeSession()

		time.Sleep(timeout)
	}
}

func (e *Exchange) analyzeSession() {
	e.activeOrderMtx.Lock()
	defer e.activeOrderMtx.Unlock()

	less := func(i, j int) bool {
		return e.activeSellOrderSlice[i].Price < e.activeSellOrderSlice[j].Price
	}
	sort.Slice(e.activeSellOrderSlice, less)

	lessForBuy := func(i, j int) bool {
		return e.activeBuyOrderSlice[i].Price < e.activeBuyOrderSlice[j].Price
	}
	sort.Slice(e.activeBuyOrderSlice, lessForBuy)

	// fmt.Println(e.activeSellOrderSlice)
	// fmt.Println(e.activeBuyOrderSlice)

	for {
		if len(e.activeSellOrderSlice) == 0 || len(e.activeBuyOrderSlice) == 0 {
			return
		}

		sellOrder := e.activeSellOrderSlice[0]
		buyOrder := e.activeBuyOrderSlice[len(e.activeBuyOrderSlice)-1]

		lowestSellPrice := sellOrder.Price
		highestBuyPrice := buyOrder.Price // last element

		if lowestSellPrice <= highestBuyPrice {
			var deal *deal

			// delete both orders
			if sellOrder.Amount == buyOrder.Amount {
				e.activeSellOrderSlice = e.activeSellOrderSlice[1:]
				e.activeBuyOrderSlice = e.activeBuyOrderSlice[:len(e.activeBuyOrderSlice)-1]

				deal = newDeal(sellOrder.Amount, math.Max(lowestSellPrice, highestBuyPrice))
				e.dealChan <- deal
				continue
			}

			// delete buy order
			if sellOrder.Amount > buyOrder.Amount {
				e.activeBuyOrderSlice = e.activeBuyOrderSlice[:len(e.activeBuyOrderSlice)-1]

				e.activeSellOrderSlice[0].Amount -= buyOrder.Amount
				deal := newDeal(buyOrder.Amount, math.Max(lowestSellPrice, highestBuyPrice))
				e.dealChan <- deal
				continue
			}

			// delete sell order
			if sellOrder.Amount < buyOrder.Amount {
				e.activeSellOrderSlice = e.activeSellOrderSlice[1:]

				e.activeBuyOrderSlice[len(e.activeBuyOrderSlice)-1].Amount -= sellOrder.Amount
				deal := newDeal(sellOrder.Amount, math.Max(lowestSellPrice, highestBuyPrice))
				e.dealChan <- deal
				continue
			}
		}

		return
	}
}
