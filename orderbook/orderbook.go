package orderbook

import (
	"errors"

	"github.com/bradenaw/juniper/container/xlist"
)

var (
	orderNotFound = errors.New("Order not found")
)

type Orderbook struct {
	symbol        string
	baseCCY       string
	quoteCCY      string
	lastFillPrice uint64
	ordersByID    map[uint64]*xlist.Node[order]
	bidSide       side
	askSide       side
}

func NewOrderbook(baseCCY string, quoteCCY string) *Orderbook {
	return &Orderbook{
		symbol:     baseCCY + "-" + quoteCCY,
		baseCCY:    baseCCY,
		quoteCCY:   quoteCCY,
		ordersByID: make(map[uint64]*xlist.Node[order], 10_000),
		bidSide:    *newSide(Bid),
		askSide:    *newSide(Ask),
	}
}

func (ob *Orderbook) AddOrder(orderID uint64, side OrderSide, price uint64, qty uint64) error {
	o := newOrder(orderID, side, price, qty)

	var fillPrice uint64
	var err error

	addOrder := func(o *xlist.Node[order]) {
		ob.ordersByID[o.Value.orderID] = o
	}

	if side == Bid {
		fillPrice, err = ob.bidSide.executeOrder(o, &ob.askSide, addOrder)
	} else {
		fillPrice, err = ob.askSide.executeOrder(o, &ob.bidSide, addOrder)
	}

	if fillPrice != 0 {
		ob.lastFillPrice = fillPrice
	}

	return err
}

func (ob *Orderbook) CancelOrder(orderID uint64) error {
	if orderNode, ok := ob.ordersByID[orderID]; ok {
		if orderNode.Value.side == Bid {
			ob.bidSide.cancelOrder(orderNode)
		} else {
			ob.askSide.cancelOrder(orderNode)
		}

		delete(ob.ordersByID, orderID)
		return nil
	}

	return orderNotFound
}

func (ob *Orderbook) AmendOrder(orderID uint64, newQty uint64) error {
	if orderNode, ok := ob.ordersByID[orderID]; ok {
		if orderNode.Value.side == Bid {
			return ob.bidSide.amendOrderQty(orderNode, newQty)
		} else {
			return ob.askSide.amendOrderQty(orderNode, newQty)
		}
	}

	return orderNotFound
}
