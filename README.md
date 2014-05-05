A fast limit order book implementation in Go.

Inspired by voyager's winning entry in the 2011 quantcup.org competition

```Bash
go test -bench=BenchmarkOrderBook
PASS
BenchmarkOrderBook	0 trade executions
63 trade executions
7426 trade executions
741709 trade executions
3704291 trade executions
5000000	       639 ns/op
```


