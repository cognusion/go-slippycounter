package slippycounter

import (
	"time"
)

// timeNumber is an internal tuple of an int64 and a Time
type timeNumber struct {
	Num  int64
	When time.Time
}

// newTimeNumberNow takes an int64 and returns a pointer to a timeNumber,
// with the Time initialized to time.Now()
func newTimeNumberNow(num int64) *timeNumber {
	return &timeNumber{
		Num:  num,
		When: time.Now(),
	}
}

// SlippyCounter is an additive-only counter that removes additions
// after they age out (slip). This is useful for "how many in the last X"
// questions. This is fully goro-safe, but not particularly recommended for:
// very short slip durations (<1 second);
// very high-rate Adds (thousands(?) per second);
// where having absolute certainty the Count is precise is at a given time is required.
type SlippyCounter struct {
	//opChan is where int64s are passed to be processed
	// currently values < 1 are ignored
	opChan chan int64
	// closeChan is a global closer
	closeChan chan struct{}
	// count is the current value of the SlippyCounter
	count int64
	// log is the ordered list of operations that have resulted in count
	log []*timeNumber
	// timeSlip is the amount of time between slips
	timeSlip time.Duration
	// closed is a quick state checker
	closed bool
}

// NewSlippyCounter takes a slip duration and returns an initialized SlippyCounter
func NewSlippyCounter(slip time.Duration) *SlippyCounter {
	sc := &SlippyCounter{
		timeSlip:  slip,
		opChan:    make(chan int64),
		closeChan: make(chan struct{}),
	}

	go sc.slipper()

	return sc
}

// Add increments an unclosed SlippyCounter by num
func (s *SlippyCounter) Add(num int) {
	if s.closed {
		// circuit breaker
		return
	}
	s.opChan <- int64(num)

}

// Count returns the last known value of a SlippyCounter.
// If the SlippyCounter has been Close()d, it will perpetually
// return the last value. i.e. it may not be 0
func (s *SlippyCounter) Count() int64 {
	return s.count
}

// Close closes a SlippyCounter from slipping, or allowing new
// Add()s. Maybe Freeze() would be more appropriate.
func (s *SlippyCounter) Close() {
	if !s.closed {
		s.closed = true
		close(s.closeChan)
	}
}

// slipper is the internal handler for Add(), Close(),
// and slip operations
func (s *SlippyCounter) slipper() {
	// Ticker is to fire of the slip operation
	var t *time.Ticker
	if s.timeSlip.Seconds() > 0 {
		// Valid ticker
		t = time.NewTicker(s.timeSlip)
		defer t.Stop()
	} else {
		// We don't want a ticker, but we have to create one
		// or the select{} will die violently
		t = time.NewTicker(1 * time.Second)
		t.Stop()
	}

SLIPLOOP:
	for {
		select {
		case num := <-s.opChan:
			// new value to add
			if num < 1 {
				// Don't waste my time
				continue
			}
			s.count += num
			s.log = append(s.log, newTimeNumberNow(num))
		case <-t.C:
			// timer tick
			if len(s.log) < 1 {
				// Don't bother
				continue
			}
			s.slip(s.timeSlip) // synchronous, so we don't have to lock

		case <-s.closeChan:
			// clean up
			break SLIPLOOP
		}
	}
}

func (s *SlippyCounter) slip(dur time.Duration) {
	then := time.Now().Add(dur * -1) // add negative time to subtract, because awkward
	last := 0
	lastCount := s.count

	for n, tn := range s.log {
		if tn != nil {
			if tn.When.Before(then) {
				// We'll want to cull this one
				last = n + 1
				lastCount -= tn.Num
			}
		}
	}

	if last > 0 {
		// Something to cull
		newLogSize := len(s.log) - last
		if newLogSize > 0 {
			// We still have items in the log
			newlog := make([]*timeNumber, newLogSize)
			s.log = append(newlog, s.log[last:]...)
		} else {
			// New log is empty
			s.log = []*timeNumber{}
		}

		s.count = lastCount
	}
}
