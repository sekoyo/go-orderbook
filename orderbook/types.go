package orderbook

type OrderSide uint64

const (
	Bid OrderSide = iota
	Ask
)
