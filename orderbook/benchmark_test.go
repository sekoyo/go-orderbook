package orderbook

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type message struct {
	eventType int
	orderID   uint64
	size      uint64
	price     uint64
	side      OrderSide
}

// func TestAddMarket2(t *testing.T) {
func BenchmarkOrderbook(b *testing.B) {
	// assert := assert.New(t)

	path, err := filepath.Abs("../test-data/AAPL_2012-06-21_34200000_57600000_message_10.csv")
	fmt.Println(path)
	check(err)
	data, err := os.ReadFile(path)

	check(err)
	reader := bytes.NewReader(data)

	r := csv.NewReader(reader)

	// messages, err := r.ReadAll()
	// check(err)

	var messages []message

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// eventType, orderID, size, price, direction := msg[1], msg[2], msg[3], msg[4], msg[5]
		eventType, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}
		orderID, err := strconv.Atoi(record[2])
		if err != nil {
			continue
		}
		size, err := strconv.Atoi(record[3])
		if err != nil {
			continue
		}
		price, err := strconv.Atoi(record[4])
		if err != nil {
			continue
		}
		direction := record[5]

		var side OrderSide
		if direction == "1" {
			side = Bid
		} else {
			side = Ask
		}

		msg := message{
			eventType: eventType,
			orderID:   uint64(orderID),
			size:      uint64(size),
			price:     uint64(price),
			side:      side,
		}

		messages = append(messages, msg)
	}

	ob := NewOrderbook("BTC", "USD")

	fmt.Println("messages len:", len(messages))

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, msg := range messages {
			if msg.eventType == 1 {
				ob.AddOrder(msg.orderID, msg.side, msg.price, msg.size)
			} else if msg.eventType == 2 {
				// TODO: is 2 an amend?
				ob.AmendOrder(msg.orderID, msg.size)
				// if err == nil {
				// fmt.Println("Amend error!", err)
				// }
			} else if msg.eventType == 3 {
				ob.CancelOrder(msg.orderID)
				// if err == nil {
				// 	fmt.Println("Cancel error!", err)")
				// }
			}
		}
	}
}
