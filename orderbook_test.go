// gorderbook_test
package gorderbook

import (
	_ "fmt"
	"testing"
)

func testOrderBook(t *testing.T) {
	cb := func(traderBuy Name, traderSell Name, price Price, size Size) {
		if traderBuy != "Kevin" {
			t.Error("expected buyer Kevin")
		}
		if traderSell != "Tom" {
			t.Error("expected seller Tom")
		}
		if price != 200 {
			t.Error("expected price 200")
		}
		if size != 100 {
			t.Error("expected size 100")
		}
	}
	book := NewOrderBook("BTCUSD", 0, 10000, cb)
	(&book).LimitOrder(Buy, 100, 593, "Kevin")
	(&book).LimitOrder(Sell, 100, 200, "Tom")
}

func TestOrderBook2(t *testing.T) {
	// See "Trading and Exchanges" by Harris, p126 (Continuous Trading Example)
	book := NewOrderBook("BTCUSD", 0, 10000, LogExecute)
	(&book).LimitOrder(Buy, 3, 200, "Bea")
	(&book).LimitOrder(Sell, 2, 201, "Sam")
	(&book).LimitOrder(Buy, 2, 200, "Ben")
	(&book).LimitOrder(Sell, 1, 198, "Sol")
	(&book).LimitOrder(Sell, 5, 202, "Stu")
	// Market order not suported so we simulate with 2 limit orders
	// at market for Bif
	(&book).LimitOrder(Buy, 2, 201, "Bif")
	(&book).LimitOrder(Buy, 2, 202, "Bif")
	(&book).LimitOrder(Buy, 2, 201, "Bob")
	(&book).LimitOrder(Sell, 6, 200, "Sue")
	(&book).LimitOrder(Buy, 7, 198, "Bud")
}
