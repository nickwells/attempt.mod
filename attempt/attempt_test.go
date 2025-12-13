package attempt_test

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/nickwells/attempt.mod/attempt"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func makeFunc(succeedAt uint64, err error) attempt.Func {
	return func(t uint64) error {
		if t < succeedAt {
			return err
		}

		return nil
	}
}

func TestAttempt(t *testing.T) {
	delayUnit := time.Millisecond * 50
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		count    uint64
		f        attempt.Func
		w        attempt.Waiter
		expCount uint64
		expDur   time.Duration
	}{
		{
			ID:       testhelper.MkID("bad count"),
			count:    0,
			f:        makeFunc(0, nil),
			w:        attempt.NoDelay{},
			expCount: 0,
			ExpErr:   testhelper.MkExpErr(attempt.BadAttemptsErr),
		},
		{
			ID:       testhelper.MkID("succeed at first attempt"),
			count:    1,
			f:        makeFunc(1, errors.New("error")),
			w:        attempt.NoDelay{},
			expCount: 1,
		},
		{
			ID:       testhelper.MkID("succeed at nth attempt"),
			count:    9,
			f:        makeFunc(3, errors.New("error")),
			w:        attempt.NoDelay{},
			expCount: 3,
		},
		{
			ID:       testhelper.MkID("fail"),
			count:    2,
			f:        makeFunc(3, errors.New("error")),
			w:        attempt.NoDelay{},
			expCount: 2,
			ExpErr:   testhelper.MkExpErr("error"),
		},
		{
			ID:       testhelper.MkID("with FixedDelay"),
			count:    11,
			f:        makeFunc(10, errors.New("error")),
			w:        attempt.NewFixedDelay(delayUnit),
			expCount: 10,
			expDur:   9 * delayUnit,
		},
		{
			ID:    testhelper.MkID("with DblDelay"),
			count: 11,
			f:     makeFunc(10, errors.New("error")),
			w: attempt.NewDblDelay(
				delayUnit,
				5*delayUnit),
			expCount: 10,
			expDur:   delayUnit * ((1 + 2 + 4) + ((9 - 3) * 5)),
		},
	}

	for _, tc := range testCases {
		synctest.Test(t, func(t *testing.T) {
			start := time.Now()
			a, err := attempt.Times(tc.count, tc.f, tc.w)
			end := time.Now()

			testhelper.DiffInt(t, tc.IDStr(), "trials", a, tc.expCount)

			testhelper.CheckExpErr(t, err, tc)

			if tc.expDur != 0 {
				testhelper.DiffTime(t, tc.IDStr(), "end time",
					start.Add(tc.expDur), end)
			}
		})
	}
}

func BenchmarkTimes(b *testing.B) {
	for b.Loop() {
		_, _ = attempt.Times(
			1001,
			makeFunc(1000, errors.New("bad")),
			attempt.NoDelay{})
	}
}
