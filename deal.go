package main

import "fmt"

type deal struct {
	// amountUSD float64
	AmountBTC float64
	Price     float64
}

func newDeal(amountBTC float64, price float64, ) *deal {
	return &deal{
		// amountUSD: amountUSD,
		AmountBTC: amountBTC,
		Price:     price,
	}
}

func (deal *deal) String() string {
	tmpl := `
	amountBTC %v
	price     %v
	`
	return fmt.Sprintf(tmpl, deal.AmountBTC, deal.Price)
}