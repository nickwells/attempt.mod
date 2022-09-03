/*
Package attempt offers funcs that let you try to run a function up to a
maximum number of times, repeating if the function returns an error. You can
also specify the interval between attempts by passing a Waiter which will
calculate the next interval between attempts.
*/
package attempt
