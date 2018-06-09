package main

import (
	"math/rand"
	"sync"
	"time"
	"fmt"
)

type priceEstimator interface {
	GetPrice() float64
}

type staticPriceEstimator struct {
	price float64
}

func newStaticPriceEstimator(price float64) priceEstimator {
	return &staticPriceEstimator{
		price: price,
	}
}

func (e *staticPriceEstimator) GetPrice() float64 {
	return e.price
}

type dynamicPriceEstimator struct {
	price float64
	priceMtx *sync.Mutex
}

func newDynamicPriceEstimator(initialPrice float64) priceEstimator {
	e :=  &dynamicPriceEstimator{
		price: initialPrice,
		priceMtx: &sync.Mutex{},
	}
	go e.activate()
	return e
}

// NOTE: must be run as goroutine
func (e *dynamicPriceEstimator) activate() priceEstimator {
	for {
		time.Sleep(time.Second * 1)

		e.priceMtx.Lock()
		increase := rand.Intn(2) == 0
		maxPercent := float64(3)
		percent := rand.Float64() * maxPercent

		if increase {
			e.price = e.price * (1.0 + percent / 100)
		} else {
			e.price = e.price * (1.0 - percent / 100)
		}
		fmt.Println(e.price)
		e.priceMtx.Unlock()
	}
}

func (e *dynamicPriceEstimator) GetPrice() float64 {
	e.priceMtx.Lock()
	defer e.priceMtx.Unlock()
	return e.price
}