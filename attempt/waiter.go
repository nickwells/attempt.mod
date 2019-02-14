package attempt

import "time"

// A Waiter is a source of delay between the calls to the function being tried.
// It's expected that it will call time.Sleep to introduce a delay between
// repeated calls to the function
type Waiter interface {
	Wait()
}

// NoDelay represents a Waiter that waits for no time
type NoDelay struct{}

// Wait returns immediately without waiting
func (w NoDelay) Wait() {
}

// DblDelay represents a Waiter that will wait for an initial period and
// then double the interval every time until it is greater than a maximum value
// after which time it will wait for the maximum time every time.
type DblDelay struct {
	firstDelay, maxDelay, delay time.Duration
}

const (
	// DfltDblDelayFirst is the initial delay that will be used if no valid
	// value is given
	DfltDblDelayFirst = time.Second
	// DfltDblDelayMax is the maximum delay that will be used if no valid
	// value is given
	DfltDblDelayMax = 5 * time.Minute
)

// NewDblDelay returns a pointer to a DblDelay struct that has the initial
// and maximum waits set to the given parameters. If either of the values is
// less than or equal to zero then it will be set to the default value when
// it is first used
func NewDblDelay(first, max time.Duration) *DblDelay {
	return &DblDelay{
		firstDelay: first,
		maxDelay:   max,
	}
}

// Wait waits for the given number of microseconds and then doubles
// that value and checks to see if it is greater than the maximum
// value and if so it sets it to the maximum value
func (w *DblDelay) Wait() {
	if w.delay == 0 {
		if w.firstDelay <= 0 {
			w.firstDelay = DfltDblDelayFirst
		}
		if w.maxDelay <= 0 {
			w.maxDelay = DfltDblDelayMax
		}
		w.delay = w.firstDelay
	}
	if w.delay > w.maxDelay {
		w.delay = w.maxDelay
	}

	time.Sleep(w.delay)
	w.delay *= 2
}

// FixedDelay represents a Waiter that will wait for a fixed period each time
type FixedDelay struct {
	delay time.Duration
}

// NewFixedDelay returns a FixedDelay struct that has the fixed delay set to
// the given parameter. The delay is measured in microseconds
func NewFixedDelay(d time.Duration) *FixedDelay {
	return &FixedDelay{
		delay: d,
	}
}

// Wait waits for the given number of microseconds
func (w FixedDelay) Wait() {
	if w.delay == 0 {
		return
	}

	time.Sleep(w.delay)
}
