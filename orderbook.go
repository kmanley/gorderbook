// a fast limit order book implementation inspired by the
// 2011 quantcup winner "voyager". See http://www.quantcup.org/

package gorderbook

import (
	"fmt"
)

type Side int8
type Size int32
type Price int32
type OrderID int32

const (
	Buy Side = iota
	Sell
)

type Order struct {
	side    Side    // buy or sell
	size    Size    // quantity
	price   Price   // price in ticks
	trader  string  // trader name
	orderID OrderID // internal order ID
}

type Orders []Order

type OrderBook struct {
	name        string          // name of the book, e.g. "USDBTC"
	orderID     OrderID         // current order ID
	maxPrice    Price           // max price in ticks supported by book
	callback    ExecuteCallback // trade execution callback
	pricePoints []Orders        // standing orders at each price level
	bidMax      Price           // current max bid
	askMin      Price           // current min ask
}

// signature for trade execution callback
type ExecuteCallback func(string, string, Price, Size)

// default trade execution function which simply logs to stdout
func LogExecute(traderBuy string, traderSell string, price Price, size Size) {
	fmt.Println("EXECUTE", traderBuy, "BUY", traderSell, size, "@", price)
}

// Creates a new order book
func NewOrderBook(name string, startOrderID OrderID, maxPrice Price,
	callback ExecuteCallback) *OrderBook {
	ob := new(OrderBook)
	ob.name = name
	ob.orderID = startOrderID
	ob.maxPrice = maxPrice
	ob.callback = callback
	ob.pricePoints = make([]Orders, maxPrice)
	ob.bidMax = 0
	ob.askMin = maxPrice + 1
	return ob
}

// Dumps a string representation of the order book to stdout. Useful
// for debugging
func (ob *OrderBook) Dump() {
	for i := ob.maxPrice - 1; i >= 0; i-- {
		entries := ob.pricePoints[i]
		if len(entries) > 0 {
			fmt.Printf("%4d: %v\n", i, entries)
		}
	}
}

/*
Inserts a new limit order into the order book. If the order
can be matched, one or more calls to the execution callback will
be made synchronously (one for each fill). If the order can't
be completely filled, it will be queued in the order book. In either
case this function returns an order ID that can be used to cancel
the order.*/
func (ob *OrderBook) LimitOrder(side Side, size Size, price Price, trader string) OrderID {
	return ob.LimitOrderEx(side, size, price, trader, ob.callback)
}

// Same as LimitOrder, but allows passing in a custom trade execution callback.
// Useful for unit testing
func (ob *OrderBook) LimitOrderEx(side Side, size Size, price Price,
	trader string, callback ExecuteCallback) OrderID {
	ob.orderID += 1
	if side == Buy {
		// look for outstanding sell orders that cross with the incoming order
		for price >= ob.askMin {
			for len(ob.pricePoints[ob.askMin]) > 0 {
				entries := ob.pricePoints[ob.askMin]
				entry := &entries[0]
				if entry.size < size {
					// the waiting entry's size is less than this buyer's size,
					// so the waiting entry is completely filled
					callback(trader, entry.trader, price, entry.size)
					size -= entry.size
					ob.pricePoints[ob.askMin] = entries[1:]
				} else {
					// the waiting entry's size is greater or equal than the buyer's size
					// so the buyer's order is completely filled
					callback(trader, entry.trader, price, size)
					if entry.size > size {
						entry.size -= size
					} else {
						ob.pricePoints[ob.askMin] = entries[1:]
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
		ob.pricePoints[price] = append(ob.pricePoints[price], Order{side, size, price, trader, ob.orderID})
		if ob.bidMax < price {
			ob.bidMax = price
		}
		return ob.orderID
	} else {
		// new order is a Sell
		// look for outstanding buy orders that cross with the incoming order
		for price <= ob.bidMax {
			for len(ob.pricePoints[ob.bidMax]) > 0 {
				entries := ob.pricePoints[ob.bidMax]
				entry := &entries[0]
				if entry.size < size {
					callback(entry.trader, trader, price, entry.size)
					size -= entry.size
					ob.pricePoints[ob.bidMax] = entries[1:]
				} else {
					callback(entry.trader, trader, price, size)
					if entry.size > size {
						entry.size -= size
					} else {
						// entry.size == size
						ob.pricePoints[ob.bidMax] = entries[1:]
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
		ob.pricePoints[price] = append(ob.pricePoints[price], Order{side, size, price, trader, ob.orderID})
		if ob.askMin > price {
			ob.askMin = price
		}
		return ob.orderID
	}
	return -1
}

/*
TODO:
 - implement market order type
 - implement order cancellation
*/
