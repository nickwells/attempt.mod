package attempt

import "errors"

// BadAttemptsErr is the error that will be returned if you pass a zero count
// to Times
const BadAttemptsErr = "the number of attempts must be greater than zero"

// Func is the signature of the function to be attempted. It will be called
// until it returns nil or the maximum number of trials has been reached
type Func func(trial uint64) error

// Times will make count attempts to call f. If f returns nil then Times
// will return the number of attempts made and nil otherwise it will keep
// trying and will return the last error returned by f when all the trials are
// up. After each failed attempt it will call w.Wait()
func Times(count uint64, f Func, w Waiter) (uint64, error) {
	if count == 0 {
		return 0, errors.New(BadAttemptsErr)
	}

	return attemptImpl(count, f, w)
}

// Forever will call f until it returns nil. After each failed attempt it
// will call w.Wait(). Note that the count of attempts passed to f can overflow
// and wrap to zero
func Forever(f Func, w Waiter) (uint64, error) {
	return attemptImpl(0, f, w)
}

func attemptImpl(count uint64, f Func, w Waiter) (uint64, error) {
	var i uint64

	var err error

	for {
		i++

		err = f(i)
		if err == nil {
			break
		}

		if count > 0 && i >= count {
			break
		}

		w.Wait()
	}

	return i, err
}
