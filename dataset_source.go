package main

import "math/rand"

const (
	minDeviation = 2
	maxDeviation = 50
)

type datasetSource struct {
	priceEstimator priceEstimator
	amountGenerator amountGenerator
}

func newDatasetSource(priceEstimator priceEstimator, amountGenerator amountGenerator) *datasetSource {
	return &datasetSource{
		priceEstimator: priceEstimator,
		amountGenerator: amountGenerator,
	}
}

func (s *datasetSource) getNextSellOrder() *order {
	price := s.priceEstimator.GetPrice() * (1.0 + getRandomDeviation())
	return newOrder(sellOrderType, price, s.amountGenerator.AmountBTC())
}

func (s *datasetSource) getNextFilteredSellOrder() *order {
	var (
		bestOrder     *order
		numIterations = 50
	)

	for i := 0; i < numIterations; i++ {
		order := s.getNextSellOrder()

		if bestOrder == nil || bestOrder.Price > order.Price {
			bestOrder = order
		}
	}
	return bestOrder
}

func (s *datasetSource) getNextBuyOrder() *order {
	price := s.priceEstimator.GetPrice() * (1.0 - getRandomDeviation())
	return newOrder(buyOrderType, price, s.amountGenerator.AmountBTC())
}

func (s *datasetSource) getNextFilteredBuyOrder() *order {
	var (
		bestOrder     *order
		numIterations = 50
	)

	for i := 0; i < numIterations; i++ {
		order := s.getNextBuyOrder()

		if bestOrder == nil || bestOrder.Price < order.Price {
			bestOrder = order
		}
	}
	return bestOrder
}

func getRandomDeviationInPercent() float64 {
	return rand.Float64()*(maxDeviation-minDeviation) + minDeviation
}

func getRandomDeviation() float64 {
	return getRandomDeviationInPercent() / 100
}
