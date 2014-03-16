/* Questions
- any way to enforce access to struct members to go through accessor function?

*/
//package orderbook
package main

import (
	"container/list"
	"fmt"
	"time"
)

type Side int8
type Size int32
type Price int32
type OrderID int32
type Name string

//type Deque list.List
type Deque struct {
	list.List
}
type ExecuteCallback func(Name, Name, Price, Size)

const (
	Buy Side = iota
	Sell
)

type Order struct {
	side    Side  // Buy or Sell
	size    Size  // quantity
	price   Price // price in ticks
	trader  Name
	orderID OrderID
}

type OrderBook struct {
	name        Name
	orderID     OrderID
	maxPrice    Price
	callback    ExecuteCallback
	pricePoints []Deque
	bidMax      Price
	askMin      Price
}

func LogExecute(traderBuy Name, traderSell Name, price Price, size Size) {
	fmt.Println("EXECUTE", traderBuy, "BUY", traderSell, size, "@", price)
}

func New(name Name, startOrderID OrderID, maxPrice Price, callback ExecuteCallback) *OrderBook {
	ob := new(OrderBook)
	ob.name = name
	ob.orderID = startOrderID
	ob.maxPrice = maxPrice
	ob.callback = callback
	ob.pricePoints = make([]Deque, maxPrice)
	ob.bidMax = 0
	ob.askMin = maxPrice + 1
	return ob
}

/*
Inserts a new limit order into the order book. If the order 
can be matched, one or more calls to the execution callback will
be made synchronously (one for each fill). If the order can't
be completely filled, it will be queued in the order book. In either
case this function returns an order ID that can be used to cancel 
the order.*/
func (ob *OrderBook) LimitOrder(side Side, size Size, price Price, trader Name) OrderID {
	ob.orderID += 1
	if side == Buy {
		// look for outstanding sell orders that cross with the incoming order
		for price >= ob.askMin {
			entries := ob.pricePoints[ob.askMin]
			for entries.Len() > 0 {
				entry := (entries.Front().Value).(Order)
				if entry.size < size {
					// the waiting entry's size is less than this buyer's size, 
					// so the waiting entry is completely filled
					ob.callback(trader, entry.trader, price, entry.size)
					size -= entry.size
					entries.Remove(entries.Front())
				} else {
					// the waiting entry's size is greater or equal than the buyer's size
					// so the buyer's order is completely filled
					ob.callback(trader, entry.trader, price, size)
					if entry.size > size {
						entry.size -= size
					} else {
						entries.Remove(entries.Front())
					}
					ob.orderID += 1
					return ob.orderID
				}
			}
			// we have exhausted all orders at the ask_min price point. Move on
			// to the next price level
			ob.askMin += 1
		}
		// if we get here then there is some qty we can't fill, so enqueue the remaining size
		ob.orderID += 1
		ob.pricePoints[price].PushBack(Order{side, size, price, trader, ob.orderID})
		if ob.bidMax < price {
			ob.bidMax = price
		}
		return ob.orderID
	} else {
		// new order is a Sell
		// look for outstanding buy orders that cross with the incoming order
		for price <= ob.bidMax {
			entries := ob.pricePoints[ob.bidMax]
			for entries.Len() > 0 {
				entry := (entries.Front().Value).(Order)
				if entry.size < size {
					ob.callback(entry.trader, trader, price, entry.size)
					size -= entry.size
					entries.Remove(entries.Front())
				} else {
					ob.callback(entry.trader, trader, price, size)
					if entry.size > size {
						entry.size -= size
					} else {
						// entry.size == size
						entries.Remove(entries.Front())
					}
					ob.orderID += 1
					return ob.orderID
				}
			}
			// we have exhausted all orders at the ask_min price point. Move on
			// to the next price level
			ob.bidMax -= 1
		}
		// if we get here then there is some qty we can't fill, so enqueue the order
		ob.orderID += 1
		ob.pricePoints[price].PushBack(Order{side, size, price, trader, ob.orderID})
		if ob.askMin > price {
			ob.askMin = price
		}
		return ob.orderID
	}
	return -1
}

func main() {
	ob := New("GBPUSD", 0, 10000, LogExecute)
	start := time.Now()
	//fmt.Println(
	(&ob).LimitOrder(Buy, 100, 593, "Kevin") //)
	fmt.Println("got here")
	//fmt.Println(
	(&ob).LimitOrder(Sell, 100, 200, "Tom") //)
	elapsed := time.Since(start)
	fmt.Println("elapsed time:", elapsed)
	/*
		fmt.Println(ob.name)
		ob.pricePoints[0].PushBack(10)
		ob.pricePoints[1].PushBack(100)
		fmt.Println(ob.pricePoints[0])
		l := ob.pricePoints[1]
		for e := l.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value)
		}
	*/
}
