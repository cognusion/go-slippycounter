package slippycounter

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCounterZeroSlip(t *testing.T) {

	sc := NewSlippyCounter(0)
	defer sc.Close()

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})
	})
}

func TestCounterIdle(t *testing.T) {

	sc := NewSlippyCounter(500 * time.Millisecond)
	defer sc.Close()

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("After adding a bunch of stuff", func() {
			for c := 0; c < 10; c++ {
				time.Sleep(50 * time.Millisecond)
				sc.Add(c)
			}

			Convey("After Time Passes", func() {
				time.Sleep(3 * time.Second)

				Convey("The Value Should still be zero", func() {
					So(sc.Count(), ShouldEqual, 0)
				})

			})
		})
	})
}

func TestCounterAddition(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	defer sc.Close()
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1)", func() {
			sc.Add(1)
			Convey("The Value Should be Greater than it was", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 1)
			})
		})

		Convey("When Add(1) Again", func() {
			sc.Add(1)

			Convey("The Value Should be Greater than it was", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 2)
			})
		})

		Convey("After Time Passes", func() {
			time.Sleep(4 * time.Second)

			Convey("The Value Should be Less than it was", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 0)
			})

		})
	})
}

func TestCounterClosingAdd(t *testing.T) {

	sc := NewSlippyCounter(2 * time.Second)
	sd := 5 * time.Millisecond

	Convey("Given new Slippy Counter", t, func() {

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("After the Counter is Closed", func() {
			sc.Close()

			Convey("When Add(1)", func() {
				sc.Add(1)
				Convey("The Value Should be the same", func() {
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

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1)", func() {
			sc.Add(1)
			Convey("The Value Should be Greater than it was", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 1)
			})
		})

		Convey("After the Counter is Closed", func() {
			sc.Close()

			Convey("After Time Passes", func() {
				time.Sleep(3 * time.Second)

				Convey("The Value Should be the same", func() {
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

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(-1)", func() {
			sc.Add(-1)

			Convey("The Value Should be the same", func() {
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

		Convey("The Value Should be zero", func() {
			So(sc.Count(), ShouldEqual, 0)
		})

		Convey("When Add(1)", func() {
			sc.Add(1)
			Convey("The Value Should be Greater than it was", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 1)
			})

		})

		Convey("After Short Time Passes", func() {
			time.Sleep(1 * time.Millisecond)

			Convey("When Add(1)", func() {
				sc.Add(1)
				Convey("The Value Should be Greater than it was", func() {
					time.Sleep(sd)
					So(sc.Count(), ShouldEqual, 2)
				})

			})

			Convey("When Add(1) Again", func() {
				sc.Add(1)
				Convey("The Value Should be Greater than it was", func() {
					time.Sleep(sd)
					So(sc.Count(), ShouldEqual, 3)
				})

			})

		})

		Convey("After Time Passes", func() {
			time.Sleep(4 * time.Second)

			Convey("The Value Should be Less than it was", func() {
				time.Sleep(sd)
				So(sc.Count(), ShouldEqual, 0)
			})

		})
	})
}
