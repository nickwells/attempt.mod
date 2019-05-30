[![GoDoc](https://godoc.org/github.com/nickwells/attempt.mod?status.png)](https://godoc.org/github.com/nickwells/attempt.mod)

# attempt
Some funcs and types for making multiple attempts with varying waits between.

There are two funcs provided: 

  * `Times` - this will try calling the supplied function up to a maximum
    number of times
  * `Forever` - this will keep trying forever.

Each of these takes a function to be attempted which should return an error
if the attempt failed. If an error is returned they will wait and try
again. You supply a `Waiter` which decides how long to wait for between
attempts.

## Waiter
The Waiter is an interface - there is a single method: `Wait()`. This package
provides some useful `Waiters`:
  * `NoDelay` - this will not wait at all but will return immediately
  * `DblDelay` - this will wait for an initial duration and this delay will
    double each time until a maximum wait duration is reached
  * `FixedDelay` - this waits for the same duration each time. The duration
    can be 0 in which case this behaves like `NoDelay`
