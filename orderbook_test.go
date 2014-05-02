// gorderbook_test
package main

import (
	"fmt"
	ob "orderbook"
	"testing"
)

func TestOrderBook(b *testing.B) {
	cb := func(traderBuy ob.Name, traderSell ob.Name, price ob.Price, size ob.Size) {
		fmt.PrintLn("hi")
	}
	book := ob.NewOrderBook("GBPUSD", 0, 10000, LogExecute)
	(&book).LimitOrder(Buy, 100, 593, "Kevin")
	(&book).LimitOrder(Sell, 100, 200, "Tom")

}
