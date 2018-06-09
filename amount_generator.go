package main

type amountGenerator interface {
	AmountUSD() float64
	AmountBTC() float64
}

type staticAmountGenerator struct {
	amountUSD float64
	amountBTC float64
}

func newStaticAmountGenerator(amountUSD float64, amountBTC float64) amountGenerator {
	return &staticAmountGenerator{
		amountUSD: amountUSD,
		amountBTC: amountBTC,
	}
}

func (g *staticAmountGenerator) AmountUSD() float64 {
	return g.amountUSD
}

func (g *staticAmountGenerator) AmountBTC() float64 {
	return g.amountBTC
}