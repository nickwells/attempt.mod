package attempt_test

import (
	"errors"
	"fmt"
	"github.com/nickwells/attempt.mod/attempt"
	"github.com/nickwells/mathutil.mod/mathutil"
	"testing"
	"time"
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
	testCases := []struct {
		testName    string
		count       uint64
		f           attempt.Func
		w           attempt.Waiter
		expCount    uint64
		errExpected bool
		expErr      string
		expDur      time.Duration
	}{
		{
			testName:    "bad count",
			count:       0,
			f:           makeFunc(0, nil),
			w:           attempt.NoDelay{},
			expCount:    0,
			errExpected: true,
			expErr:      attempt.BadAttemptsErr,
		},
		{
			testName: "succeed at first attempt",
			count:    1,
			f:        makeFunc(1, errors.New("error")),
			w:        attempt.NoDelay{},
			expCount: 1,
		},
		{
			testName: "succeed at nth attempt",
			count:    9,
			f:        makeFunc(3, errors.New("error")),
			w:        attempt.NoDelay{},
			expCount: 3,
		},
		{
			testName:    "fail",
			count:       2,
			f:           makeFunc(3, errors.New("error")),
			w:           attempt.NoDelay{},
			expCount:    2,
			errExpected: true,
			expErr:      "error",
		},
		{
			testName: "with FixedDelay",
			count:    101,
			f:        makeFunc(100, errors.New("error")),
			w:        attempt.NewFixedDelay(10 * time.Millisecond),
			expCount: 100,
			expDur:   99 * 10 * time.Millisecond,
		},
		{
			testName: "with DblDelay",
			count:    101,
			f:        makeFunc(100, errors.New("error")),
			w: attempt.NewDblDelay(
				time.Millisecond,
				5*time.Millisecond),
			expCount: 100,
			expDur:   time.Millisecond * ((1 + 2 + 4) + ((99 - 3) * 5)),
		},
	}

	for i, tc := range testCases {
		start := time.Now()
		a, err := attempt.Times(tc.count, tc.f, tc.w)
		end := time.Now()
		testID := fmt.Sprintf("test %d: %s :", i, tc.testName)

		if a != tc.expCount {
			t.Errorf("%s expected to stop after %d tries actual: %d",
				testID, tc.expCount, a)
		}

		if err == nil && tc.errExpected {
			t.Errorf("%s expected err: %s but none returned", testID, tc.expErr)
		} else if err != nil {
			if !tc.errExpected {
				t.Errorf("%s no error was expected but we saw: %s",
					testID, err)
			} else if err.Error() != tc.expErr {
				t.Errorf("%s expected err: %s but we saw: %s",
					testID, tc.expErr, err)
			}
		}

		if tc.expDur != 0 {
			dur := end.Sub(start)
			pct := 5.0
			if !mathutil.WithinNPercent(float64(dur), float64(tc.expDur), pct) {
				diff := (dur - tc.expDur)
				t.Logf("%s:\n", testID)
				t.Logf("\t:   actual duration: %6d ms\n",
					time.Duration(dur.Nanoseconds())/time.Millisecond)
				t.Logf("\t: expected duration: %6d ms\n",
					time.Duration(tc.expDur.Nanoseconds())/time.Millisecond)
				t.Logf("\t:        difference: %9d µs (%5.1f%%)\n",
					diff/time.Microsecond,
					100.0*float64(diff)/float64(tc.expDur))
				t.Errorf(
					"\t: difference is not within %.1f%% of expected value",
					pct)
			}
		}
	}
}

func BenchmarkTimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		attempt.Times(
			1001,
			makeFunc(1000, errors.New("bad")),
			attempt.NoDelay{})
	}
}
