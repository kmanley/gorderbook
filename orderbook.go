//package orderbook
package gorderbook

import (
	_ "bytes"
	"fmt"
)

type Side int8
type Size int32
type Price int32
type OrderID int32
type Name string

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

type Orders []Order

type OrderBook struct {
	name        Name
	orderID     OrderID
	maxPrice    Price
	callback    ExecuteCallback
	pricePoints []Orders
	bidMax      Price
	askMin      Price
}

func LogExecute(traderBuy Name, traderSell Name, price Price, size Size) {
	fmt.Println("EXECUTE", traderBuy, "BUY", traderSell, size, "@", price)
}

func NewOrderBook(name Name, startOrderID OrderID, maxPrice Price, callback ExecuteCallback) *OrderBook {
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
func (ob *OrderBook) LimitOrder(side Side, size Size, price Price, trader Name) OrderID {
	fmt.Printf("%+v %+v %+v %+v\n", side, size, price, trader)
	ob.orderID += 1
	if side == Buy {
		// look for outstanding sell orders that cross with the incoming order
		for price >= ob.askMin {
			for len(ob.pricePoints[ob.askMin]) > 0 {
				//fmt.Printf("entries at %d: %+v\n", ob.askMin, entries)
				//fmt.Println("----------------------------------------------")
				entries := ob.pricePoints[ob.askMin]
				entry := entries[0]
				if entry.size < size {
					// the waiting entry's size is less than this buyer's size,
					// so the waiting entry is completely filled
					ob.callback(trader, entry.trader, price, entry.size)
					size -= entry.size
					//fmt.Println("len before: ", len(entries))
					ob.pricePoints[ob.askMin] = entries[1:]
					//fmt.Println("len after: ", len(entries))
				} else {
					// the waiting entry's size is greater or equal than the buyer's size
					// so the buyer's order is completely filled
					ob.callback(trader, entry.trader, price, size)
					if entry.size > size {
						entry.size -= size
					} else {
						//fmt.Println("len before: ", len(entries))
						ob.pricePoints[ob.askMin] = entries[1:]
						//fmt.Println("len after: ", len(entries))
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
				//fmt.Printf("entries at %d: %+v\n", ob.bidMax, entries)
				//fmt.Println("----------------------------------------------")
				entries := ob.pricePoints[ob.bidMax]
				entry := entries[0]
				if entry.size < size {
					ob.callback(entry.trader, trader, price, entry.size)
					size -= entry.size
					//fmt.Println("len before: ", len(entries))
					ob.pricePoints[ob.bidMax] = entries[1:]
					//fmt.Println("len after: ", len(entries))
				} else {
					ob.callback(entry.trader, trader, price, size)
					if entry.size > size {
						entry.size -= size
					} else {
						// entry.size == size
						//fmt.Println("len before: ", len(entries))
						ob.pricePoints[ob.bidMax] = entries[1:]
						//fmt.Println("len after: ", len(entries))
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
func (ob *OrderBook) renderLevel(level Level) {
	var buffer bytes.Buffer
	for

        ret = ",".join(("%s:%s(%s)" % (order.size, order.trader, order.order_id) for order in level))
        if len(ret) > maxlen:
            ret = ",".join((str(order.size) for order in level))
        if len(ret) > maxlen:
            ret = "%d orders (total size %d)" % (len(level), sum((order.size for order in level)))
            assert len(ret) <= maxlen
        return ret
		}
*/
