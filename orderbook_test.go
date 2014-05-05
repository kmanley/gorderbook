// gorderbook_test
package gorderbook

import (
	"fmt"
	"math/rand"
	"testing"
)

func assertEqual(t *testing.T, lhs interface{}, rhs interface{}) {
	if lhs != rhs {
		t.Error(lhs, " != ", rhs)
	}
}

func TestOrderBook(t *testing.T) {
	cb := func(traderBuy string, traderSell string, price Price, size Size) {
		assertEqual(t, traderBuy, "Kevin")
		assertEqual(t, traderSell, "Tom")
		assertEqual(t, price, Price(200))
		assertEqual(t, size, Size(100))
	}
	book := NewOrderBook("BTCUSD", 0, 10000, cb)
	(&book).LimitOrder(Buy, 100, 593, "Kevin")
	(&book).LimitOrder(Sell, 100, 200, "Tom")
}

func TestOrderBook2(t *testing.T) {
	// See "Trading and Exchanges" by Harris, p126 (Continuous Trading Example)
	//func
	book := NewOrderBook("BTCUSD", 0, 10000, LogExecute)
	(&book).LimitOrderEx(Buy, 3, 200, "Bea", func(_ string, _ string, _ Price, _ Size) {
		t.Error("shouldn't have been called")
	})
	(&book).LimitOrderEx(Sell, 2, 201, "Sam", func(_ string, _ string, _ Price, _ Size) {
		t.Error("shouldn't have been called")
	})
	(&book).LimitOrderEx(Buy, 2, 200, "Ben", func(_ string, _ string, _ Price, _ Size) {
		t.Error("shouldn't have been called")
	})
	(&book).LimitOrderEx(Sell, 1, 198, "Sol", func(buyer string, seller string, price Price, size Size) {
		assertEqual(t, buyer, "Bea")
		assertEqual(t, seller, "Sol")
		assertEqual(t, price, Price(198))
		assertEqual(t, size, Size(1))
	})
	(&book).LimitOrderEx(Sell, 5, 202, "Stu", func(_ string, _ string, _ Price, _ Size) {
		t.Error("shouldn't have been called")
	})
	// Market order not suported so we simulate with 2 limit orders
	// at market for Bif
	(&book).LimitOrderEx(Buy, 2, 201, "Bif", func(buyer string, seller string, price Price, size Size) {
		assertEqual(t, buyer, "Bif")
		assertEqual(t, seller, "Sam")
		assertEqual(t, price, Price(201))
		assertEqual(t, size, Size(2))
	})
	(&book).LimitOrderEx(Buy, 2, 202, "Bif", func(buyer string, seller string, price Price, size Size) {
		assertEqual(t, buyer, "Bif")
		assertEqual(t, seller, "Stu")
		assertEqual(t, price, Price(202))
		assertEqual(t, size, Size(2))
	})
	(&book).LimitOrderEx(Buy, 2, 201, "Bob", func(_ string, _ string, _ Price, _ Size) {
		t.Error("shouldn't have been called")
	})
	ctr := 0
	(&book).LimitOrderEx(Sell, 6, 200, "Sue", func(buyer string, seller string, price Price, size Size) {
		ctr += 1
		if ctr == 1 {
			assertEqual(t, buyer, "Bob")
			assertEqual(t, seller, "Sue")
			assertEqual(t, price, Price(200))
			assertEqual(t, size, Size(2))
		} else if ctr == 2 {
			assertEqual(t, buyer, "Bea")
			assertEqual(t, seller, "Sue")
			assertEqual(t, price, Price(200))
			assertEqual(t, size, Size(2))
		} else if ctr == 3 {
			assertEqual(t, buyer, "Ben")
			assertEqual(t, seller, "Sue")
			assertEqual(t, price, Price(200))
			assertEqual(t, size, Size(2))
		} else {
			t.Error("too many trade executions!")
		}

	})
	(&book).LimitOrderEx(Buy, 7, 198, "Bud", func(_ string, _ string, _ Price, _ Size) {
		t.Error("shouldn't have been called")
	})
}

// go test -bench=BenchmarkOrderBook
// BenchmarkOrderBook	 5000000	       635 ns/op
// ok  	github.com/kmanley/gorderbook	3.835s
func BenchmarkOrderBook(b *testing.B) {
	minSize := 1
	maxSize := 20
	minPrice := 8000
	maxPrice := 9500
	ctr := 0
	book := NewOrderBook("BTCUSD", 0, 10000, func(string, string, Price, Size) { ctr += 1 })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		price := rand.Intn(maxPrice-minPrice) + minPrice
		size := rand.Intn(maxSize-minSize) + minSize
		var side Side
		if rand.Intn(1000) >= 500 {
			side = Buy
		} else {
			side = Sell
		}
		(&book).LimitOrder(side, Size(size), Price(price), "Trader")
	}
	fmt.Println(ctr, "trade executions")
}
