package orderbook

import (
	"github.com/bradenaw/juniper/container/xlist"
	"github.com/bradenaw/juniper/xmath"
)

type level struct {
	totalQty uint64
	orders   xlist.List[order]
}

// Expects check has been made to see if order is
// fillable against this level.
func (l *level) fill(levelPrice uint64, o *order) uint64 {
	var fillSize uint64 = 0

	for ok := true; ok; ok = l.orders.Len() > 0 {
		restingOrderNode := l.orders.Front()
		restingOrder := &restingOrderNode.Value
		qtyToFill := xmath.Min(o.qtyLeft(), restingOrder.qtyLeft())

		restingOrder.fill(levelPrice, qtyToFill)
		o.fill(levelPrice, qtyToFill)

		l.totalQty -= qtyToFill
		fillSize += qtyToFill

		if restingOrder.isFilled() {
			l.orders.Remove(restingOrderNode)
		}

		if o.isFilled() {
			break
		}
	}

	return fillSize
}
