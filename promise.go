package promise

// SuccessHandler is the function prototype for promise listeners that
// receive the results of a successful delivery of the promise
type SuccessHandler func(result interface{})

// CatchHandler is the function prototype for promise listeners that
// receive the error from an unsuccessful promise delivery
type CatchHandler func(err error)

// CanceledHandler is the function prototype for promise listeners that
// receive a callback when promise delivery is canceled
type CanceledHandler func()

// AlwaysHandler is the function prototype for promise listeners that
// receive a callback regardless of the result of the promise deliver
type AlwaysHandler func(promise Controller)

// Factory is a function prototype that returns a Promise
type Factory func() Promise

// FactoryWithResult is used to pass the result of a promise to a function
// that creates another promise
type FactoryWithResult func(result interface{}) Promise

// Promise is the interface for Promise delivery
type Promise interface {
	// Success registers a callback on successful delivery of the promise
	Success(handler SuccessHandler) Promise

	// Catch registers a callback on a failed delivery of the promise
	Catch(handler CatchHandler) Promise

	// Canceled registers a callback for the case where the promise delivery
	// is canceled
	Canceled(handler CanceledHandler) Promise

	// Always registers a callback when the promise is delivered or canceled
	Always(handler AlwaysHandler) Promise

	// Allows a wait on promise delivery via a channel
	//
	//  Notes
	//		Blocks until the promise is delivered
	//
	//    Equivalent to:
	//      p.Always(func (p Controller) {
	//        myChan <- p
	//      })
	//
	//		return <-myChan
	//
	Wait(chan Controller) Promise

	// Use a channel as a signal when the promise is delivered without
	// blocking
	//
	//  Notes
	//    Equivalent to:
	//      p.Always(func (p Controller) {
	//        myChan <- p
	//      })
	//
	//		return p
	//
	Signal(waitChan chan Controller) Promise

	// Chain a Promise to the successful delivery of this Promise
	Then(promise Promise) Promise

	// Chain a Promise (created via Factory) to the successful delivery of
	// this Promise
	Thenf(factory Factory) Promise

	// ThenWithResult chains the result of a successful promise to another
	// promise
	ThenWithResult(factory FactoryWithResult) Promise

	// Chain a list of Promises to the successful delivery of this Promise
	ThenAll(promises ...Promise) Promise

	// Chain a list of Promises (created via Factory) to the successful
	// delivery of this Promise
	ThenAllf(factory func() []Promise) Promise

	// Chain a promise to successful delivery of any one from a list of Promises after
	// successful delivery of this Promise
	ThenAny(promises ...Promise) Promise

	// Chain a promise to successful delivery of any one from a list of Promises
	// after (created via Factory) successful delivery of this Promise
	ThenAnyf(factories func() []Promise) Promise
}
