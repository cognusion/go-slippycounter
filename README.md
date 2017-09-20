# go-slippycounter


[![GoDoc](https://godoc.org/github.com/cognusion/go-slippycounter?status.svg)](https://godoc.org/github.com/cognusion/go-slippycounter)

SlippyCounter is an additive-only, eventually-consistent counter that removes additions after they age out (slip). 
This is very useful for long-running "how many in the last X" questions. 

This is fully goro-safe, but not particularly recommended for: 
* Very short slip durations (<1 second)
* ~~Very high-rate Adds (thousands(?) per second)~~ This is actually fine now
* Where having absolute certainty the Count() is precise is at a given time is required

Count() may be more or less than what you expect.

## Benchmarks

* BenchmarkAdd1 repeatedly adds to an existing slipless SlippyCounter
* BenchmarkSlipZero repeatedly slips an empty SlippyCounter
* BenchmarkSlip1k repeatedly slips a SlippyCounter with 1000 slippy items (worst case)
* BenchmarkSlip10k repeatedly slips a SlippyCounter with 10000 slippy items (worst case)

```BASH
BenchmarkAdd1-8       	 3000000	       451 ns/op	      78 B/op	       1 allocs/op
BenchmarkSlipZero-8   	30000000	        53.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkSlip1k-8     	  300000	      4887 ns/op	       0 B/op	       0 allocs/op
BenchmarkSlip10k-8    	   30000	     49958 ns/op	       0 B/op	       0 allocs/op
```