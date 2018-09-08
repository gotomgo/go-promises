package promise

import "fmt"

// ErrPromiseCanceled is used as the error result when a Promise is canceled
var ErrPromiseCanceled = fmt.Errorf("The promise delivery was canceled")

// Controller is an interface for controlling the state / result of
// the Promise
type Controller interface {
	Promise

	// Result returns the successful result of the delivery or nil
	//
	//  Notes
	//    A successful result may BE nil, so if you need to disambiguate,
	//    use IsSuccess() which returns true IF the promise has been delivered
	//		and it was not an error
	//
	Result() interface{}

	// Error returns the error for a failed delivery or nil if the
	// result is not a failure
	//
	//  Notes
	//    See IsFailed(), IsError(), and Fail(error)
	//
	//    In the case of a canceled promise delivery, this method will
	//    return ErrPromiseCanceled
	//
	Error() error

	// RawResult returns the underlying result / error
	//
	//  Notes
	//    RawResult is useful for transferring the result of a promise to
	//    another promise
	//
	RawResult() interface{}

	// Succeed delivers the promise with a value of true
	Succeed() Controller

	// SucceedWithResult delivers the promise successfully with the specified
	// result
	SucceedWithResult(result interface{}) Controller

	// DeliverWithPromise delivers the promise based on the result of a
	// different Promise (Controller)
	DeliverWithPromise(promise Controller) Controller

	// Deliver delivers the promise and based on the type of the result,
	// determines the success or failure
	//
	//  Notes
	//    if result is of type error, then Fail(result.(error)), otherwise
	//    SucceedWithResult(result)
	//
	Deliver(result interface{}) Controller

	// Fail fails the deliver of the promise with an error
	Fail(err error) Controller

	// Cancel cancels the promise
	//
	//  Notes
	//    The value of Error() will return ErrPromiseCanceled for a canceled
	//    Promise
	Cancel() Controller

	// IsPending determins if the promise is still pending delivery
	IsPending() bool

	// IsDelivered determines if the promise has been delivered
	IsDelivered() bool

	// IsSuccess determines if the promise has been successfully delivered
	IsSuccess() bool

	// IsFailed determines if the promise has been delivered with an error
	IsFailed() bool

	// IsError determines if the promise has been delivered with an error
	//
	//  Notes
	//    IsError is an alias for IsFailed
	//
	IsError() bool

	// IsCanceled determines if the promise delivery has been canceled
	IsCanceled() bool
}
