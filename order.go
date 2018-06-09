package main

import "fmt"

type orderType uint8

const (
	sellOrderType orderType = iota
	buyOrderType
)

func (orderType orderType) String() string {
	switch orderType {
	case sellOrderType:
		return "Sell order"
	case buyOrderType:
		return "Buy order"
	}
	return "<Unknown type>"
}

type order struct {
	OrderType orderType
	Price     float64
	Amount    float64
}

func newOrder(orderType orderType, price float64, amount float64) *order {
	return &order{
		OrderType: orderType,
		Price:     price,
		Amount:    amount,
	}
}

func (order *order) String() string {
	tmpl := `
	orderType: %v
	price:     %v
	`
	return fmt.Sprintf(tmpl, order.OrderType, order.Price)
}
