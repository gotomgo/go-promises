package promise

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUndelivered(t *testing.T) {
	p := NewPromise()

	assert.False(t, p.IsDelivered())
	assert.True(t, p.IsPending())
	assert.False(t, p.IsFailed())
	assert.False(t, p.IsError())
	assert.False(t, p.IsCanceled())

	// when not delivered, this should return false
	assert.False(t, p.IsSuccess())

	assert.Nil(t, p.Error())
	assert.Nil(t, p.Result())
	assert.Nil(t, p.RawResult())

}

func TestFailPromise(t *testing.T) {
	p := NewPromise()

	testErr := fmt.Errorf("Testing promise.Fail(err)")

	// let's fail the promise
	p2 := p.Fail(testErr)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.True(t, p.IsFailed())
	assert.True(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.False(t, p.IsSuccess())

	assert.Equal(t, testErr, p.Error())
	assert.Nil(t, p.Result())
	assert.Equal(t, testErr, p.RawResult())

	//	assert.Panics(t, func() { p.Fail(testErr) })
}

func TestCancelPromise(t *testing.T) {
	p := NewPromise()

	// let's cancel the promise
	p2 := p.Cancel()
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.True(t, p.IsFailed())
	assert.True(t, p.IsError())
	assert.True(t, p.IsCanceled())
	assert.False(t, p.IsSuccess())

	assert.Equal(t, ErrPromiseCanceled, p.Error())
	assert.Nil(t, p.Result())
	assert.Equal(t, ErrPromiseCanceled, p.RawResult())

	// assert.Panics(t, func() { p.Cancel() })
}

func TestSucceedPromise(t *testing.T) {
	p := NewPromise()

	// let's cancel the promise
	p2 := p.Succeed()
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.False(t, p.IsFailed())
	assert.False(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.True(t, p.IsSuccess())

	assert.Equal(t, nil, p.Error())
	assert.Equal(t, true, p.Result())
	assert.Equal(t, true, p.RawResult())

	// assert.Panics(t, func() { p.Succeed() })
}

func TestSucceedWithResultPromise(t *testing.T) {
	p := NewPromise()

	// let's cancel the promise
	p2 := p.SucceedWithResult(12)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.False(t, p.IsFailed())
	assert.False(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.True(t, p.IsSuccess())

	assert.Equal(t, nil, p.Error())
	assert.Equal(t, 12, p.Result())
	assert.Equal(t, 12, p.RawResult())

	// assert.Panics(t, func() { p.Succeed() })
}

func TestDeliverError(t *testing.T) {
	p := NewPromise()

	testErr := fmt.Errorf("Testing promise.Fail(err)")

	// let's fail the promise
	p2 := p.Deliver(testErr)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.True(t, p.IsFailed())
	assert.True(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.False(t, p.IsSuccess())

	assert.Equal(t, testErr, p.Error())
	assert.Nil(t, p.Result())
	assert.Equal(t, testErr, p.RawResult())

	// assert.Panics(t, func() { p.Fail(testErr) })
}

func TestDeliverCancel(t *testing.T) {
	p := NewPromise()

	// let's cancel the promise
	p2 := p.Deliver(ErrPromiseCanceled)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.True(t, p.IsFailed())
	assert.True(t, p.IsError())
	assert.True(t, p.IsCanceled())
	assert.False(t, p.IsSuccess())

	assert.Equal(t, ErrPromiseCanceled, p.Error())
	assert.Nil(t, p.Result())
	assert.Equal(t, ErrPromiseCanceled, p.RawResult())

	// assert.Panics(t, func() { p.Cancel() })
}

func TestDeliverSuccess(t *testing.T) {
	p := NewPromise()

	p2 := p.Deliver(12)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.False(t, p.IsFailed())
	assert.False(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.True(t, p.IsSuccess())

	assert.Equal(t, nil, p.Error())
	assert.Equal(t, 12, p.Result())
	assert.Equal(t, 12, p.RawResult())

	// assert.Panics(t, func() { p.Succeed() })
}

func TestDeliverNil(t *testing.T) {
	p := NewPromise()

	p2 := p.Deliver(nil)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.False(t, p.IsFailed())
	assert.False(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.True(t, p.IsSuccess())

	assert.Equal(t, nil, p.Error())
	assert.Equal(t, nil, p.Result())
	assert.Equal(t, nil, p.RawResult())

	// assert.Panics(t, func() { p.Succeed() })
}

func TestDeliverWithPromise(t *testing.T) {
	p := NewPromise()
	other := NewPromise()
	other.Deliver(12)

	p2 := p.DeliverWithPromise(other)
	assert.Equal(t, p, p2)

	assert.True(t, p.IsDelivered())
	assert.False(t, p.IsPending())
	assert.False(t, p.IsFailed())
	assert.False(t, p.IsError())
	assert.False(t, p.IsCanceled())
	assert.True(t, p.IsSuccess())

	assert.Equal(t, nil, p.Error())
	assert.Equal(t, 12, p.Result())
	assert.Equal(t, 12, p.RawResult())

	// assert.Panics(t, func() { p.Succeed() })
}

func TestDeliverWithPendingPromise(t *testing.T) {
	p := NewPromise()
	other := NewPromise()

	assert.Panics(t, func() { p.DeliverWithPromise(other) })
}

func TestPostSuccessNotify(t *testing.T) {
	p := NewPromise()
	p.Deliver(12)

	var onSuccess int

	p.Success(func(result interface{}) {
		onSuccess++
		assert.Equal(t, 12, result)
	})

	assert.Equal(t, 1, onSuccess)
}

func TestPostFailNotify(t *testing.T) {
	p := NewPromise()

	testErr := fmt.Errorf("Test fail notify")
	p.Deliver(testErr)

	var onCatch int

	p.Catch(func(err error) {
		onCatch++
		assert.Equal(t, testErr, err)
	})

	assert.Equal(t, 1, onCatch)
}

func TestPostCancelNotify(t *testing.T) {
	p := NewPromise()

	p.Cancel()

	var onCanceled int

	p.Canceled(func() {
		onCanceled++
	})

	assert.Equal(t, 1, onCanceled)
}

func TestPostAlwaysNotify(t *testing.T) {
	p := NewPromise()

	p.Succeed()

	var onAlways int

	p.Always(func(p2 Controller) {
		onAlways++
		assert.Equal(t, p, p2)
	})

	assert.Equal(t, 1, onAlways)
}

func TestPreSuccessNotify(t *testing.T) {
	p := NewPromise()

	var onSuccess int
	var onAlways int
	var onFailed int
	var onCanceled int

	p.Success(func(result interface{}) {
		onSuccess++
		assert.Equal(t, 12, result)
	}).Always(func(p2 Controller) {
		onAlways++
		assert.Equal(t, p, p2)
	}).Catch(func(err error) {
		onFailed++
	}).Canceled(func() {
		onCanceled++
	})

	p.Deliver(12)

	assert.Equal(t, 1, onSuccess)
	assert.Equal(t, 1, onAlways)
	assert.Equal(t, 0, onFailed)
	assert.Equal(t, 0, onCanceled)
}

func TestPreFailNotify(t *testing.T) {
	p := NewPromise()

	var onSuccess int
	var onAlways int
	var onFailed int
	var onCanceled int

	testErr := fmt.Errorf("Testing pre fail notifications")

	p.Success(func(result interface{}) {
		onSuccess++
	}).Always(func(p2 Controller) {
		onAlways++
		assert.Equal(t, p, p2)
	}).Catch(func(err error) {
		onFailed++
		assert.Equal(t, testErr, err)
	}).Canceled(func() {
		onCanceled++
	})

	p.Fail(testErr)

	assert.Equal(t, 0, onSuccess)
	assert.Equal(t, 1, onAlways)
	assert.Equal(t, 1, onFailed)
	assert.Equal(t, 0, onCanceled)
}

func TestPreCancelNotify(t *testing.T) {
	p := NewPromise()

	var onSuccess int
	var onAlways int
	var onFailed int
	var onCanceled int

	p.Success(func(result interface{}) {
		onSuccess++
	}).Always(func(p2 Controller) {
		onAlways++
		assert.Equal(t, p, p2)
	}).Catch(func(err error) {
		onFailed++
		assert.Equal(t, ErrPromiseCanceled, err)
	}).Canceled(func() {
		onCanceled++
	})

	p.Cancel()

	assert.Equal(t, 0, onSuccess)
	assert.Equal(t, 1, onAlways)
	assert.Equal(t, 1, onFailed)
	assert.Equal(t, 1, onCanceled)
}

func TestPostWaitNotify(t *testing.T) {
	p := NewPromise()

	p.Succeed()

	go func(p Controller) {
		var onWait int

		myChan := make(chan Controller)

		p.Wait(myChan)

		select {
		case p2 := <-myChan:
			onWait++
			assert.Equal(t, p, p2)
		case <-time.After(1 * time.Second):
		}

		assert.Equal(t, 1, onWait)
	}(p)
}

func TestPreWaitNotify(t *testing.T) {
	p := NewPromise()

	defer p.Succeed()

	go func(p Controller) {
		var onWait int

		myChan := make(chan Controller)

		p.Wait(myChan)

		select {
		case p2 := <-myChan:
			onWait++
			assert.Equal(t, p, p2)
		case <-time.After(1 * time.Second):
		}

		assert.Equal(t, 1, onWait)
	}(p)
}

func TestThenPromise(t *testing.T) {
	p := NewPromise()
	p.Succeed()

	var onThen int
	var myChan = make(chan Controller)

	thenP := p.Then(deferredPromiseFunc()).Success(func(result interface{}) {
		onThen++
	}).Wait(myChan)

	assert.True(t, thenP.(Controller).IsDelivered())
	assert.Equal(t, 1, onThen)
}

func TestThenfPromise(t *testing.T) {
	p := NewPromise()
	p.Succeed()

	var onThen int
	var myChan = make(chan Controller)

	thenP := p.Thenf(deferredPromiseFunc).Success(func(result interface{}) {
		onThen++
	}).Wait(myChan)

	assert.True(t, thenP.(Controller).IsDelivered())
	assert.Equal(t, 1, onThen)
}

func deferredPromiseFunc() Promise {
	p := NewPromise()

	go func(p Controller) {
		select {
		case <-time.After(500 * time.Millisecond):
			p.Succeed()
		}
	}(p)

	return p
}

func TestThenfFailedPromise(t *testing.T) {
	p := NewPromise()

	testErr := fmt.Errorf("Testing pre fail notifications")

	p.Fail(testErr)

	var onThen int
	var myChan = make(chan Controller, 1)

	thenP := p.Thenf(deferredPromiseFunc).Success(func(result interface{}) {
		onThen++
	}).Wait(myChan)

	assert.True(t, thenP.(Controller).IsError())
	assert.Equal(t, 0, onThen)
}

func TestThenAll(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onSuccess int

	NewPromise().Succeed().ThenAll(p1, p2).Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 1, onSuccess)
}

func TestThenAllFail(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onCancel int

	NewPromise().Cancel().ThenAll(p1, p2).Canceled(func() {
		onCancel++
	})

	assert.Equal(t, 1, onCancel)
}

func TestThenfAll(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onSuccess int

	NewPromise().Succeed().ThenAllf(func() []Promise { return []Promise{p1, p2} }).Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 1, onSuccess)
}

func TestThenAllfFail(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onCancel int

	NewPromise().Cancel().ThenAllf(func() []Promise { return []Promise{p1, p2} }).Canceled(func() {
		onCancel++
	})

	assert.Equal(t, 1, onCancel)
}

func TestThenAllPostFail(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Cancel()

	var onCancel int

	NewPromise().Succeed().ThenAll(p1, p2).Canceled(func() {
		onCancel++
	})

	assert.Equal(t, 1, onCancel)
}

func TestThenAllEmpty(t *testing.T) {
	var onSuccess int

	NewPromise().Succeed().ThenAll().Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 1, onSuccess)
}

func TestThenAny(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onSuccess int

	NewPromise().Succeed().ThenAny(p1, p2).Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 1, onSuccess)
}

func TestThenAnyFail(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onSuccess int

	NewPromise().Cancel().ThenAny(p1, p2).Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 0, onSuccess)
}

func TestThenAnyf(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onSuccess int

	NewPromise().Succeed().ThenAnyf(func() []Promise { return []Promise{p1, p2} }).Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 1, onSuccess)
}

func TestThenAnyfFail(t *testing.T) {
	p1 := NewPromise().Succeed()
	p2 := NewPromise().Succeed()

	var onSuccess int

	NewPromise().Cancel().ThenAnyf(func() []Promise { return []Promise{p1, p2} }).Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 0, onSuccess)
}

func TestThenAnyEmpty(t *testing.T) {
	var onSuccess int

	NewPromise().Succeed().ThenAny().Success(func(result interface{}) {
		onSuccess++
	})

	assert.Equal(t, 1, onSuccess)
}

func TestPostSignalNotify(t *testing.T) {
	p := NewPromise()

	p.Succeed()

	go func(p Controller) {
		var onWait int

		myChan := make(chan Controller)

		p.Signal(myChan)

		select {
		case p2 := <-myChan:
			onWait++
			assert.Equal(t, p, p2)
		case <-time.After(1 * time.Second):
		}

		assert.Equal(t, 1, onWait)
	}(p)
}

func TestPreSignalNotify(t *testing.T) {
	p := NewPromise()

	defer p.Succeed()

	go func(p Controller) {
		var onWait int

		myChan := make(chan Controller)

		p.Signal(myChan)

		select {
		case p2 := <-myChan:
			onWait++
			assert.Equal(t, p, p2)
		case <-time.After(1 * time.Second):
		}

		assert.Equal(t, 1, onWait)
	}(p)

	// force a sleep so the GO routine can get going prior to Succeed()
	// Note that the time required is non-deterministic, but 800 seems to work
	time.Sleep(800)
}

func TestBadSuccessHandler(t *testing.T) {
	p := NewPromise()

	var onSuccess int
	p.Success(func(result interface{}) {
		// cause a panic
		panic(fmt.Errorf("test panic"))
	}).Success(func(result interface{}) {
		onSuccess++
	})

	p.Succeed()

	assert.Equal(t, 1, onSuccess)
}

func TestBadCatchHandler(t *testing.T) {
	p := NewPromise()

	var onCatch int
	p.Catch(func(err error) {
		panic(fmt.Errorf("test panic"))
	}).Catch(func(err error) {
		onCatch++
	})

	p.Fail(fmt.Errorf("test"))

	assert.Equal(t, 1, onCatch)
}

func TestBadCanceledHandler(t *testing.T) {
	p := NewPromise()

	var onCanceled int
	p.Canceled(func() {
		panic(fmt.Errorf("test panic"))
	}).Canceled(func() {
		onCanceled++
	})

	p.Cancel()

	assert.Equal(t, 1, onCanceled)
}

func TestBadAlwaysHandler(t *testing.T) {
	p := NewPromise()

	var onAlways int
	p.Always(func(p2 Controller) {
		panic(fmt.Errorf("test panic"))
	}).Always(func(p2 Controller) {
		onAlways++
	})

	p.Succeed()

	assert.Equal(t, 1, onAlways)
}
