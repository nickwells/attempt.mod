package attempt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/nickwells/attempt.mod/attempt"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
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
		start := time.Now()
		a, err := attempt.Times(tc.count, tc.f, tc.w)
		end := time.Now()

		testhelper.DiffInt(t, tc.IDStr(), "trials", a, tc.expCount)

		testhelper.CheckExpErr(t, err, tc)

		if tc.expDur != 0 {
			dur := end.Sub(start)
			diff := (dur - tc.expDur)

			if diff < 0 {
				t.Log(tc.IDStr())
				t.Logf("\t:   actual duration: %6d ms\n",
					time.Duration(dur.Nanoseconds())/time.Millisecond)
				t.Logf("\t: expected duration: %6d ms\n",
					time.Duration(tc.expDur.Nanoseconds())/time.Millisecond)
				t.Error("\t: finished sooner than expected")

				continue
			}

			const pct = 6.0

			if !mathutil.WithinNPercent(float64(dur), float64(tc.expDur), pct) {
				t.Log(tc.IDStr())
				t.Logf("\t:   actual duration: %6d ms\n",
					time.Duration(dur.Nanoseconds())/time.Millisecond)
				t.Logf("\t: expected duration: %6d ms\n",
					time.Duration(tc.expDur.Nanoseconds())/time.Millisecond)
				t.Logf("\t:        difference: %9d µs (%5.1f%%)\n",
					diff/time.Microsecond,
					100.0*float64(diff)/float64(tc.expDur))
				t.Errorf("\t: difference is more than %.1f%% of expected value",
					pct)
			}
		}
	}
}

func BenchmarkTimes(b *testing.B) {
	for range b.N {
		_, _ = attempt.Times(
			1001,
			makeFunc(1000, errors.New("bad")),
			attempt.NoDelay{})
	}
}
