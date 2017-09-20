# go-slippycounter


[![GoDoc](https://godoc.org/github.com/cognusion/go-slippycounter?status.svg)](https://godoc.org/github.com/cognusion/go-slippycounter)

SlippyCounter is an additive-only, eventually-consistent counter that removes additions after they age out (slip). 
This is very useful for long-running "how many in the last X" questions. 

This is fully goro-safe, but not particularly recommended for: 
* Very short slip durations (<1 second)
* Very high-rate Adds (thousands(?) per second)
* Where having absolute certainty the Count() is precise is at a given time is required

Count() may be more or less than what you expect.