package slippycounter

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// Creates a 30 second slippy counter, belches the
// counter value every 1 second, and does some Add()
// and Sleep()
func ExampleSlippyCounter() {

	counter := NewSlippyCounter(30 * time.Second)
	defer counter.Close()

	// Goro to belch the count every second
	go func() {
		for c := 0; c < 60; c++ {
			fmt.Printf("Count: %d\n", counter.Count())
			time.Sleep(1 * time.Second)
		}
	}()

	counter.Add(1)
	counter.Add(3)
	time.Sleep(2 * time.Second)
	counter.Add(5)
	counter.Add(7)
	time.Sleep(10 * time.Second)
	counter.Add(9)

	time.Sleep(40 * time.Second)

}

func TestCounterZeroSlip(t *testing.T) {

	sc := NewSlippyCounter(0)
	defer sc.Close()

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})
	})
}

func TestCounterIdle(t *testing.T) {

	sc := NewSlippyCounter(500 * time.Millisecond)
	defer sc.Close()

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("After adding a bunch of stuff", func() {
			for c := 0; c < 10; c++ {
				time.Sleep(50 * time.Millisecond)
				sc.Add(c)
			}

			Convey("The Value should not be zero", func() {
				So(sc.Count(), ShouldNotEqual, 0)
			})

			Convey("After Time Passes", func() {
				time.Sleep(3 * time.Second)

				Convey("The Value should again be zero", func() {
					So(sc.Count(), ShouldEqual, 0)
				})

			})
		})
	})
}

func TestCounterAdditionSlipping(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	defer sc.Close()
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1)", func() {
			sc.Add(1)
			Convey("The Value should be one", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 1)
			})
		})

		Convey("When Add(1) Again", func() {
			sc.Add(1)

			Convey("The Value should be two", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 2)
			})
		})

		Convey("After Time Passes", func() {
			time.Sleep(4 * time.Second)

			Convey("The Value should be zero", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 0)
			})

		})
	})
}

func TestCounterAddition30k(t *testing.T) {

	sc := NewSlippyCounter(0)
	defer sc.Close()
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1) 30k times", func() {
			for c := 0; c < 30000; c++ {
				sc.Add(1)
			}
			Convey("The Value should be 30k", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 30000)
			})
		})
	})
}

func TestCounterClosingAdd(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("After the Counter is Closed", func() {
			sc.Close()

			Convey("When Add(1)", func() {
				sc.Add(1)
				Convey("The Value should still be zero", func() {
					time.Sleep(sd)
					So(sc.Count(), ShouldEqual, 0)
				})

			})

		})
	})
}

func TestCounterClosingLatent(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1)", func() {
			sc.Add(1)
			Convey("The Value should be one", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 1)
			})
		})

		Convey("After the Counter is Closed", func() {
			sc.Close()

			Convey("After Time Passes", func() {
				time.Sleep(3 * time.Second)

				Convey("The Value should still be one", func() {
					time.Sleep(sd)
					So(sc.Count(), ShouldEqual, 1)
				})

			})
		})
	})

}

func TestCounterSubZero(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	defer sc.Close()
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(-1)", func() {
			sc.Add(-1)

			Convey("The Value should still be zero", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 0)
			})

		})
	})

}

func TestCounterBurn(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	defer sc.Close()
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1)", func() {
			sc.Add(1)
			Convey("The Value should be one", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 1)
			})

		})

		Convey("After Short Time Passes", func() {
			time.Sleep(1 * time.Millisecond)

			Convey("When Add(1)", func() {
				sc.Add(1)
				Convey("The Value should be two", func() {
					time.Sleep(sd)
					So(sc.Count(), ShouldEqual, 2)
				})

			})

			Convey("When Add(1) Again", func() {
				sc.Add(1)
				Convey("The Value should be three", func() {
					time.Sleep(sd)
					So(sc.Count(), ShouldEqual, 3)
				})

			})

		})

		Convey("After Time Passes", func() {
			time.Sleep(4 * time.Second)

			Convey("The Value should be zero", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 0)
			})

		})
	})
}

func BenchmarkAdd1(b *testing.B) {
	sc := NewSlippyCounter(0)
	defer sc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Add(1)
	}
}

func BenchmarkSlipZero(b *testing.B) {
	sc := NewSlippyCounter(0)
	defer sc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.slip(0)
	}
}

func BenchmarkSlip1k(b *testing.B) {
	sc := NewSlippyCounter(0)
	defer sc.Close()

	for c := 0; c < 1000; c++ {
		sc.Add(1)
	}

	// cache these
	scl := sc.log
	scc := sc.count

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// replace key parts
		sc.log = scl
		sc.count = scc

		// Slip it all
		sc.slip(0)
	}
}

func BenchmarkSlip10k(b *testing.B) {
	sc := NewSlippyCounter(0)
	defer sc.Close()

	for c := 0; c < 10000; c++ {
		sc.Add(1)
	}

	// cache these
	scl := sc.log
	scc := sc.count

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// replace key parts
		sc.log = scl
		sc.count = scc

		// Slip it all
		sc.slip(0)
	}
}
