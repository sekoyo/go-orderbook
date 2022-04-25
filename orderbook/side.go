package orderbook

import (
	"errors"

	"github.com/bradenaw/juniper/container/tree"
	"github.com/bradenaw/juniper/container/xlist"
	"golang.org/x/exp/constraints"
)

type side struct {
	orderSide OrderSide
	levels    tree.Map[uint64, level]
	depth     uint64
	bestPrice uint64
}

var (
	notEnoughLiquidity = errors.New("Not enough liquidity")
	amendTooLow        = errors.New("Cannot amend lower than filled amount")
)

func newSide(orderSide OrderSide) *side {
	var levels tree.Map[uint64, level]

	if orderSide == Bid {
		levels = tree.NewMap[uint64, level](OrderedMore[uint64])
	} else {
		levels = tree.NewMap[uint64, level](OrderedLess[uint64])
	}

	return &side{
		orderSide: orderSide,
		levels:    levels,
	}
}

func (s *side) executeOrder(o *order, oppSide *side, addOrderCb func(o *xlist.Node[order])) (uint64, error) {
	// Market order needs to fully fill.
	if o.price == 0 && oppSide.depth < o.qty {
		return 0, notEnoughLiquidity
	}

	var lastFillPrice uint64 = 0
	updateBestPrice := false
	iter := oppSide.levels.Iterate()

	for {
		levelNode, ok := iter.Next()

		if !ok || !oppSide.canMatch(o.price, levelNode.Key) {
			break
		}

		level := &levelNode.Value
		levelPrice := levelNode.Key

		// Fill.
		fillSize := level.fill(levelPrice, o)
		oppSide.depth -= fillSize
		lastFillPrice = levelPrice

		// Level is now empty, remove it.
		if level.orders.Len() == 0 {
			if oppSide.bestPrice == levelPrice {
				updateBestPrice = true
			}
			oppSide.levels.Delete(levelPrice)
		}

		if o.isFilled() {
			break
		}
	}

	// If level with best price got removed, update it.
	if updateBestPrice {
		oppSide.bestPrice = oppSide.getBestPrice()
	}

	// Add to book if not filled.
	if !o.isFilled() {
		if o.price == 0 {
			panic("Market orders should have fully filled")
		}

		s.addToBook(o, addOrderCb)
	}

	return lastFillPrice, nil
}

func (s *side) addToBook(o *order, addOrderCb func(o *xlist.Node[order])) {
	qtyLeft := o.qtyLeft()
	s.depth += qtyLeft

	level := s.levels.Get(o.price)
	isNewLevel := level.totalQty == 0

	// Update level's qty.
	level.totalQty += qtyLeft

	// Point to this level for fast lookup.
	o.level = &level

	// Push order to level and save node to ordersByID.
	// Note: this copies the order.
	orderNode := level.orders.PushBack(*o)
	addOrderCb(orderNode)

	if isNewLevel {
		s.levels.Put(o.price, level)
	}

	// Update best price.
	if s.orderSide == Bid {
		if o.price > s.bestPrice {
			s.bestPrice = o.price
		}
	} else {
		if s.bestPrice == 0 || o.price < s.bestPrice {
			s.bestPrice = o.price
		}
	}
}

func (s *side) getBestPrice() uint64 {
	price, _ := s.levels.First()
	return price
}

func (s *side) canMatch(orderPrice uint64, restingPrice uint64) bool {
	if orderPrice == 0 {
		// Market orders are hoes.
		return true
	}

	if s.orderSide == Bid {
		// Sell offer must be <= buy offer.
		return orderPrice <= restingPrice
	} else {
		// Buy offer must be >= sell offer.
		return restingPrice != 0 && orderPrice >= restingPrice
	}
}

func (s *side) cancelOrder(orderNode *xlist.Node[order]) {
	o := &orderNode.Value
	o.level.totalQty -= o.qtyLeft()
	s.depth -= o.qtyLeft()

	o.level.orders.Remove(orderNode)
	if o.level.orders.Len() == 0 {
		s.levels.Delete(o.price)
		if s.bestPrice == o.price {
			s.bestPrice = s.getBestPrice()
		}
	}
}

func (s *side) amendOrderQty(orderNode *xlist.Node[order], newQty uint64) error {
	o := &orderNode.Value

	if newQty < o.qtyFilled {
		return amendTooLow
	}

	qtyDiff := newQty - o.qty

	o.qty = newQty
	o.level.totalQty += qtyDiff
	s.depth += qtyDiff

	return nil
}

// Map ordering.
func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}
func OrderedMore[T constraints.Ordered](a, b T) bool {
	return a > b
}
