package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	assert := assert.New(t)

	// Order 10 lots for 200.
	o := newOrder(0, Bid, 200, 10)

	// Fill 5 lots at 190.
	o.fill(190, 5)
	assert.Equal(uint64(5), o.qtyFilled)
	assert.Equal(float64(190), o.avgFillPrice())
	assert.Equal(false, o.isFilled())

	// Fill remaing at 200.
	o.fill(200, 5)
	assert.Equal(uint64(10), o.qtyFilled)
	assert.Equal(float64(195), o.avgFillPrice())
	assert.Equal(true, o.isFilled())
}
