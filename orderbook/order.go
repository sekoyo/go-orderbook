package orderbook

import (
	"time"
)

type order struct {
	orderID   uint64
	side      OrderSide
	price     uint64
	qty       uint64
	qtyFilled uint64
	totalCost uint64
	timestamp int64
	level     *level
}

func newOrder(orderID uint64, side OrderSide, price uint64, qty uint64) *order {
	return &order{
		orderID:   orderID,
		side:      side,
		price:     price,
		qty:       qty,
		timestamp: time.Now().Unix(),
	}
}

// Assumes only filling the qty left.
func (o *order) fill(price uint64, qty uint64) {
	o.totalCost += qty * price
	o.qtyFilled += qty
}

func (o *order) isFilled() bool {
	return o.qty-o.qtyFilled == 0
}

func (o *order) qtyLeft() uint64 {
	return o.qty - o.qtyFilled
}

func (o *order) avgFillPrice() float64 {
	if o.qtyFilled == 0 {
		return 0.0
	}
	return float64(o.totalCost) / float64(o.qtyFilled)
}
