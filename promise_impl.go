package promise

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

// promise implements Controller and Promise
type promise struct {
	// lock is used to protect use of handler arrays, and delivery
	lock             sync.Mutex
	successHandlers  []SuccessHandler
	catchHandlers    []CatchHandler
	alwaysHandlers   []AlwaysHandler
	canceledHandlers []CanceledHandler

	// the result of the promise as an atomic value
	result atomic.Value
}

var _ Controller = &promise{}

// use internally when delivery is called with value == nil
// to allow nil to be a delivered value
var nilResult = &struct{}{}

// resolved is used in cases where we want to return a successul promise
var resolved = NewPromise().Succeed()

// NewPromise creates an instance of promise which implements Controller
// (and therefore, implements Promise)
func NewPromise() Controller {
	return &promise{}
}

// IsDelivered determines if the promise has been delivered
func (p *promise) IsDelivered() bool {
	return p.result.Load() != nil
}

// IsPending determins if the promise is still pending deliver
func (p *promise) IsPending() bool {
	return !p.IsDelivered()
}

// Result returns the error for a failed delivery or nil if the
// result is not a failure
//
//  Notes
//    See IsFailed(), IsError(), and Fail(error)
//
//    In the case of a canceled promise delivery, this method will
//    return ErrPromiseCanceled
//
func (p *promise) Error() (err error) {
	res := p.result.Load()

	if res != nil {
		err, _ = res.(error)
	}

	return
}

// RawResult returns the underlying result / error
//
//  Notes
//    RawResult is useful for transferring the result of a promise to
//    another promise
//
func (p *promise) RawResult() interface{} {
	res := p.result.Load()
	if res == nilResult {
		res = nil
	}

	return res
}

// Result returns the successful result of the delivery or nil
//
//  Notes
//    A successful result may BE nil, so if you need to disambiguate,
//    use IsSuccess() which returns true IF the promise has been delivered
//		and it was not an error
//
func (p *promise) Result() interface{} {
	res := p.result.Load()

	if res != nil {
		// if the result represents an error return nil
		if _, ok := res.(error); ok {
			return nil
		}

		// if the delivery was == nil, then we use nilResult to indicate
		// and now we need to translate it back to nil
		if res == nilResult {
			return nil
		}
	}

	return res
}

// IsFailed determines if the promise has been delivered with an error
func (p *promise) IsFailed() bool {
	return p.Error() != nil
}

// IsError determines if the promise has been delivered with an error
//
//  Notes
//    IsError is an alias for IsFailed
//
func (p *promise) IsError() bool {
	return p.IsFailed()
}

// IsSuccess determines if the promise has been successfully delivered
func (p *promise) IsSuccess() bool {
	res := p.result.Load()

	// res will be nil if the promise hasnt been delivered
	if res == nil {
		return false
	}

	// success depends on whether res is / isn't an error
	_, ok := res.(error)

	return !ok
}

// IsCanceled determines if the promise delivery has been canceled
func (p *promise) IsCanceled() bool {
	return p.Error() == ErrPromiseCanceled
}

// notifySuccess invokes a SuccessHandler with panic recovery
func (p *promise) notifySuccess(handler SuccessHandler, result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("success handler panic'd: %s", r)
		}
	}()

	handler(result)
}

// notifyAlways invokes an AlwaysHandler with panic recovery
func (p *promise) notifyAlways(handler AlwaysHandler) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("always handler panic'd: %s", r)
		}
	}()

	handler(p)
}

// notifyCatch invokes a CatchHandler with panic recovery
func (p *promise) notifyCatch(handler CatchHandler, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("catch handler panic'd: %s", r)
		}
	}()

	handler(err)
}

// notifyCanceled invokes a CanceledHandler with panic recovery
func (p *promise) notifyCanceled(handler CanceledHandler) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("canceled handler panic'd: %s", r)
		}
	}()

	handler()
}

// copySuccessHandlers creates a copy of the handlers for notification
func (p *promise) copySuccessHandlers() []SuccessHandler {
	p.lock.Lock()
	defer p.lock.Unlock()

	handlers := make([]SuccessHandler, len(p.successHandlers))
	copy(handlers, p.successHandlers)

	return handlers
}

// copyCatchHandlers creates a copy of the handlers for notification
func (p *promise) copyCatchHandlers() []CatchHandler {
	p.lock.Lock()
	defer p.lock.Unlock()

	handlers := make([]CatchHandler, len(p.catchHandlers))
	copy(handlers, p.catchHandlers)

	return handlers
}

// copyAlwaysHandlers creates a copy of the handlers for notification
func (p *promise) copyAlwaysHandlers() []AlwaysHandler {
	p.lock.Lock()
	defer p.lock.Unlock()

	handlers := make([]AlwaysHandler, len(p.alwaysHandlers))
	copy(handlers, p.alwaysHandlers)

	return handlers
}

// copyCanceledHandlers creates a copy of the handlers for notification
func (p *promise) copyCanceledHandlers() []CanceledHandler {
	p.lock.Lock()
	defer p.lock.Unlock()

	handlers := make([]CanceledHandler, len(p.canceledHandlers))
	copy(handlers, p.canceledHandlers)

	return handlers
}

// notify invokes the appropriate callbacks based on the delivered result
// of the promise
//
//	Notes
//		Canceled delivery is also considered an error, so both cancel and catch
//		handlers are invoked
//
//		As the name suggests, always handlers are always invoked
//
//		Because we cannot take the lock during the notification, we must
//		make a copy of the appropriate handler arrays so they are not modified
//		while we are notifying
//
func (p *promise) notify() {
	if p.IsSuccess() {
		res := p.Result()

		handlers := p.copySuccessHandlers()
		for _, handler := range handlers {
			p.notifySuccess(handler, res)
		}
	} else {
		err := p.Error()

		// invoke the catch handlers, even if err == ErrPromiseCanceled
		handlers := p.copyCatchHandlers()
		for _, handler := range handlers {
			p.notifyCatch(handler, err)
		}

		// if canceled, invoke cancel handlers
		if err == ErrPromiseCanceled {
			handlers := p.copyCanceledHandlers()
			for _, handler := range handlers {
				p.notifyCanceled(handler)
			}
		}
	}

	handlers := p.copyAlwaysHandlers()
	for _, handler := range handlers {
		p.notifyAlways(handler)
	}
}

// deliver implements the core logic for Promise delivery
func (p *promise) deliver(result interface{}) Controller {
	var wasDelivered bool

	p.lock.Lock()
	defer func() {
		// release the lock prior to notifying
		p.lock.Unlock()

		// do we need to notify
		if wasDelivered {
			p.notify()
		}
	}()

	if !p.IsDelivered() {
		// invoke callbacks via notify()
		wasDelivered = true

		// if nil is delivered, use nilResult as a non-nil place holder
		if result == nil {
			result = nilResult
		}

		// store the delivered result
		p.result.Store(result)
	} else {
		// This would be great as a panic, but in 'all' and 'any' scenarios it
		// is difficult to prevent async code from double completing
		log.Println("Attempt to deliver promise that is already delivered")
	}

	return p
}

// Allows a wait on promise delivery via a channel
//
//  Notes
//		Blocks until the promise is delivered
//
func (p *promise) Wait(waitChan chan Controller) Promise {
	p.Always(func(p2 Controller) {
		waitChan <- p2
	})

	return <-waitChan
}

// Use a channel as a signal when the promise is delivered without
// blocking
func (p *promise) Signal(waitChan chan Controller) Promise {
	p.Always(func(p2 Controller) {
		waitChan <- p2
	})

	return p
}

// Cancel cancels the promise
//
//  Notes
//    The value of Error() will return ErrPromiseCanceled for a canceled
//    Promise
//
//		For notification purposes, Cancel is considered an error
//
func (p *promise) Cancel() Controller {
	return p.deliver(ErrPromiseCanceled)
}

// Fail fails the delivery of the promise with an error
func (p *promise) Fail(err error) Controller {
	return p.deliver(err)
}

// Succeed delivers the promise with a value of true
func (p *promise) Succeed() Controller {
	return p.deliver(true)
}

// SucceedWithResult delivers the promise successfully with the specified
// result
func (p *promise) SucceedWithResult(result interface{}) Controller {
	return p.deliver(result)
}

// DeliverWithPromise delivers the promise based on the result of a
// different Promise (Controller)
func (p *promise) DeliverWithPromise(promise Controller) Controller {
	if promise.IsPending() {
		panic(fmt.Errorf("Cannot deliver with pending promise"))
	}

	return p.deliver(promise.RawResult())
}

// Deliver delivers the promise and based on the type of the result,
// determines the success or failure
//
//  Notes
//    if result is of type error, then Fail(result.(error)), otherwise
//    SucceedWithResult(result), unless result is a Controller, then
//		equivalent to DeliverWithPromise
//
func (p *promise) Deliver(result interface{}) Controller {
	if result != nil {
		if promise, ok := result.(Controller); ok {
			return p.DeliverWithPromise(promise)
		}
	}

	return p.deliver(result)
}

// Success registers a callback on successful delivery of the promise
//
//	Notes
//		This method locks so that concurrent access cannot change
//		from !delivered to delivered while we are in this routine,
//    as that would lead to non-deterministic invokcation of the
//		callback
//
//		If the promise is already delivered when this nethod is called
//		then invocation of the callback is synchronous, otherwise it
//		is non-synchronous
//
func (p *promise) Success(handler SuccessHandler) Promise {
	var notify bool

	p.lock.Lock()
	defer func() {
		// release the lock (before calling handler)
		p.lock.Unlock()

		// do we need to directly notify?
		if notify {
			handler(p.Result())
		}
	}()

	// already delivered and successful?
	if p.IsSuccess() {
		// direct invoke
		notify = true
	} else {
		// deferred invoke
		p.successHandlers = append(p.successHandlers, handler)
	}

	return p
}

// Catch registers a callback on a failed delivery of the promise
//
//	Notes
//		This method locks so that concurrent access cannot change
//		from !delivered to delivered while we are in this routine,
//    as that would lead to non-deterministic invokcation of the
//		callback
//
//		If the promise is already delivered when this nethod is called
//		then invocation of the callback is synchronous, otherwise it
//		is non-synchronous
//
func (p *promise) Catch(handler CatchHandler) Promise {
	var notify bool

	p.lock.Lock()
	defer func() {
		// release the lock (before invoking handler)
		p.lock.Unlock()

		// is direct notify?
		if notify {
			handler(p.Error())
		}
	}()

	// is delivered and error?
	if p.IsError() {
		// direct invoke
		notify = true
	} else {
		// deferred invoke
		p.catchHandlers = append(p.catchHandlers, handler)
	}

	return p
}

// Canceled registers a callback for the case where the promise delivery
// is canceled
//
//	Notes
//		This method locks so that concurrent access cannot change
//		from !delivered to delivered while we are in this routine,
//    as that would lead to non-deterministic invokcation of the
//		callback
//
//		If the promise is already delivered when this nethod is called
//		then invocation of the callback is synchronous, otherwise it
//		is non-synchronous
//
func (p *promise) Canceled(handler CanceledHandler) Promise {
	var notify bool

	p.lock.Lock()
	defer func() {
		// release the lock (before invoking handler)
		p.lock.Unlock()

		// is direct notify?
		if notify {
			handler()
		}
	}()

	// is delivered and canceled?
	if p.IsCanceled() {
		// direct invoke
		notify = true
	} else {
		// deferred invoke
		p.canceledHandlers = append(p.canceledHandlers, handler)
	}

	return p
}

// Always registers a callback when the promise is delivered or canceled
//
//	Notes
//		This method locks so that concurrent access cannot change
//		from !delivered to delivered while we are in this routine,
//    as that would lead to non-deterministic invokcation of the
//		callback
//
//		If the promise is already delivered when this nethod is called
//		then invocation of the callback is synchronous, otherwise it
//		is non-synchronous
//
func (p *promise) Always(handler AlwaysHandler) Promise {
	var notify bool

	p.lock.Lock()
	defer func() {
		// release the lock
		p.lock.Unlock()

		// is direct notify?
		if notify {
			handler(p)
		}
	}()

	// if its delivered then direct notify
	if p.IsDelivered() {
		// direct invoke
		notify = true
	} else {
		// deferred invoke
		p.alwaysHandlers = append(p.alwaysHandlers, handler)
	}

	return p
}

// Chain a Promise to the successful delivery of this Promise
func (p *promise) Then(promise Promise) Promise {
	return p.Thenf(func() Promise { return promise })
}

// Chain a Promise (created via Factory) to the successful delivery of
// this Promise
func (p *promise) Thenf(factory Factory) Promise {
	result := NewPromise()

	p.Always(func(p2 Controller) {
		if p2.IsSuccess() {
			factory().Always(func(p3 Controller) {
				result.DeliverWithPromise(p3)
			})
		} else {
			result.DeliverWithPromise(p2)
		}
	})

	return result
}

// ThenWithResult chains the result of a successful promise to another
// promise
func (p *promise) ThenWithResult(factory FactoryWithResult) Promise {
	result := NewPromise()

	p.Always(func(p2 Controller) {
		if p2.IsSuccess() {
			factory(p2.Result()).Always(func(p3 Controller) {
				result.DeliverWithPromise(p3)
			})
		} else {
			result.DeliverWithPromise(p2)
		}
	})

	return result
}

// all is a base implementtion of ThenAll
func (p *promise) all(promises []Promise) Promise {
	// how many promises must complete?
	count := int64(len(promises))

	// none? return success
	if count == 0 {
		return resolved
	}

	// create a promise to bridge this promise and the 'all' promises
	result := NewPromise()

	for _, promise := range promises {
		// attach an always handler and based on the result do the right thing
		promise.Always(func(p2 Controller) {
			// if the promise failed, then result is failed
			if p2.IsFailed() {
				result.DeliverWithPromise(p2)
			} else {
				// once all promises complete successfully, result is successful
				if atomic.AddInt64(&count, -1) == 0 {
					result.Succeed()
				}
			}
		})

		// early-out in case the promise got delivered synchronously
		if result.IsDelivered() {
			break
		}
	}

	return result
}

// Chain a list of Promises to the successful delivery of this Promise
func (p *promise) ThenAll(promises ...Promise) Promise {
	return p.Then(p.all(promises))
}

// Chain a list of Promises (created via Factory) to the successful
// delivery of this Promise
func (p *promise) ThenAllf(factory func() []Promise) Promise {
	return p.Then(p.all(factory()))
}

// any is a base implementation of ThenAny
func (p *promise) any(promises []Promise) Promise {
	// if there are no any promises, then success
	if len(promises) == 0 {
		return resolved
	}

	// create a bridge promise between this promise and the any promises
	result := NewPromise()

	for _, promise := range promises {
		// add an always handler for each promise
		promise.Always(func(p2 Controller) {
			// deliver result based on result of promise. For Any, we only need
			// one promise to deliver, not all of them (see all([]Promise))
			result.DeliverWithPromise(p2)
		})

		// early-out in case the promise got delivered synchronously
		if result.IsDelivered() {
			break
		}
	}

	return result
}

// Chain a list of Promises to the successful delivery of this Promise
func (p *promise) ThenAny(promises ...Promise) Promise {
	return p.Then(p.any(promises))
}

// Chain a list of Promises (created via Factory) to the successful
// delivery of this Promise
func (p *promise) ThenAnyf(factory func() []Promise) Promise {
	return p.Then(p.any(factory()))
}
