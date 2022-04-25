package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMarket(t *testing.T) {
	assert := assert.New(t)

	ob := NewOrderbook("BTC", "USD")

	// Market order should fail (not enough liquidity)
	_, err := ob.AddOrder(0, Bid, 0, 100)
	assert.Equal(true, err != nil)

	// Add 2 bids with total qty 200
	ob.AddOrder(1, Bid, 4000, 100)
	ob.AddOrder(2, Bid, 4001, 100)
	assert.Equal(uint64(4001), ob.bidSide.bestPrice)

	// Market order should fail (not enough liquidity)
	_, err = ob.AddOrder(3, Bid, 0, 199)
	assert.Equal(true, err != nil)

	// Market order should wipe out bids
	mkt, err := ob.AddOrder(4, Ask, 0, 200)
	assert.Equal(true, err == nil)
	assert.Equal(true, mkt.isFilled())
	assert.Equal(uint64(0), ob.bidSide.bestPrice)
	assert.Equal(uint64(0), ob.bidSide.depth)
	assert.Equal(uint64(0), ob.askSide.bestPrice)
	assert.Equal(uint64(0), ob.askSide.depth)

}

func TestAddLimit(t *testing.T) {
	assert := assert.New(t)

	ob := NewOrderbook("BTC", "USD")
	assert.Equal(uint64(0), ob.bidSide.depth)
	assert.Equal(uint64(0), ob.askSide.depth)
	assert.Equal(uint64(0), ob.bidSide.bestPrice)
	assert.Equal(uint64(0), ob.askSide.bestPrice)
	assert.Equal(uint64(0), ob.lastFillPrice)

	// 1000 bid should rest.
	bid, _ := ob.AddOrder(0, Bid, 1000, 100)
	assert.Equal(uint64(0), ob.lastFillPrice)
	assert.Equal(uint64(1000), ob.bidSide.bestPrice)
	assert.Equal(uint64(100), ob.bidSide.depth)

	// 1100 bid should rest above.
	ob.AddOrder(1, Bid, 1100, 100)
	assert.Equal(uint64(0), ob.lastFillPrice)
	assert.Equal(uint64(1100), ob.bidSide.bestPrice)
	assert.Equal(uint64(200), ob.bidSide.depth)

	// Ask at 1200 should rest.
	ask, _ := ob.AddOrder(2, Ask, 1200, 100)
	assert.Equal(uint64(0), ob.lastFillPrice)
	assert.Equal(ob.askSide.depth, uint64(100))
	assert.Equal(ob.askSide.bestPrice, uint64(1200))
	// Nothing should have happened to the bid.
	assert.Equal(ob.bidSide.bestPrice, uint64(1100))
	assert.Equal(ob.bidSide.depth, uint64(200))

	// Ask should fill first and part of next bid.
	ask, _ = ob.AddOrder(3, Ask, 800, 120)
	assert.Equal(uint64(1000), ob.bidSide.bestPrice)
	assert.Equal(uint64(80), ob.bidSide.depth)
	assert.Equal(true, ask.isFilled())
	assert.Equal(uint64(120), ask.qtyFilled)
	assert.Equal(float64(1083.3333333333333), ask.avgFillPrice())

	// Bid should fill remaining 1200 ask and rest.
	bid, _ = ob.AddOrder(4, Bid, 1200, 120)
	assert.Equal(uint64(0), ob.askSide.bestPrice)
	assert.Equal(uint64(0), ob.askSide.depth)
	assert.Equal(uint64(100), bid.qtyFilled)
	assert.Equal(false, bid.isFilled())
	// Best bid now 1200 and 20 more to depth.
	assert.Equal(uint64(1200), ob.bidSide.bestPrice)
	assert.Equal(uint64(80+20), ob.bidSide.depth)
}

func TestCancelOrder(t *testing.T) {
	assert := assert.New(t)

	ob := NewOrderbook("BTC", "USD")

	// 1000 bid should rest.
	bid, _ := ob.AddOrder(5, Bid, 1000, 100)
	assert.Equal(uint64(0), ob.lastFillPrice)
	assert.Equal(uint64(1000), ob.bidSide.bestPrice)
	assert.Equal(uint64(100), ob.bidSide.depth)

	// Cancel bid.
	err := ob.CancelOrder(bid.orderID)
	assert.Equal(true, err == nil)
	assert.Equal(uint64(0), ob.bidSide.bestPrice)
	assert.Equal(uint64(0), ob.bidSide.depth)
}

func TestAmendOrder(t *testing.T) {
	assert := assert.New(t)

	ob := NewOrderbook("BTC", "USD")

	// 1000 bid should rest.
	bid, _ := ob.AddOrder(0, Bid, 1000, 100)
	assert.Equal(uint64(0), ob.lastFillPrice)
	assert.Equal(uint64(1000), ob.bidSide.bestPrice)
	assert.Equal(uint64(100), ob.bidSide.depth)

	// Amend bid.
	err := ob.AmendOrder(bid.orderID, 200)
	assert.Equal(true, err == nil)
	assert.Equal(uint64(200), ob.bidSide.depth)
}
