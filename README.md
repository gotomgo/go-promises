# GO Promises package

go-promises is a flexibile, easy to use implementation of promises written in GO. This version is heavily based on a C# implementation I wrote for use with the Unity 3D game engine, which has been used / tested extensively in successful commercial products. The primary difference in the implementations (beyond language) is that Unity 3D is primarily single-threaded with a co-routine model, whereas the GO implementation must support concurrent operations. For those interested in the C# version, it can be found here: https://bitbucket.org/codemules/promises

## Getting Started

1. Download and install it:

```sh
$ go get github.com/gotomgo/go-promises
```

2. Import it in your code:

```go
import promises "github.com/gotomgo/go-promises"
```

3. Review the provided examples

The examples can be found in the examples folder or at https://github.com/gotomgo/go-promises/tree/master/examples

## Basic Usage
There are two primary interfaces exposed by the promises package:
* Controller
* Promise

Controller is a super-set of Promise and is used by code that controls delivery of the promise, while
Promise is intended for consumers that react to delivery of the promise.

### Creating a Promise
The function, NewPromise, returns a Controller (which again, is a Promise) but methods that return an instance of a promise typically return Promise.

Here is a basic template for a method that creates and returns a Promise:

```go
func myFunction() promises.Promise {
  // NewPromise returns Controller, which is a superset of Promise
  p := promises.NewPromise()

  // perform work here. Synchronously or Asynchronously.
  // don't forget to deliver the Promise !!

  return p
}
```

### Using a Promise

Consumers of a Promise are primarily interested in the delivery of a promise, and if it was successfully delivered, processing the result. They are also typically interested in errors as well.

```go
func consumer() {
  myFunction().Success(func(result interface{}) {
    // do something with the result
  }).Catch(func (err error) {
    fmt.Println("Error from myFunction: ",err)
  })
}
```

In other cases you may want to use Always to perform an action or process the result.

```go
func consumer() {
  myFunction().Always(func(p promises.Controller) {
    // do something now that the promise is delivered

    someFunctionUnrelatedToResult()

    // or process the result like follows:
    if p.IsSuccess() {
      result := p.Result()
    } else {
      fmt.Println("Error from myFunction: ",p.Error())
    }
  })
```

### Delivering a Promise

The Controller interface provides a variety of methods to deliver a Promise:

```go
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
```
### Using Wait and Signal
Because GO channels are so useful as a syncrhonization mechanism, you might want to combine them with promises in some cases.

```go
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
```

### Using Then, Thenf, ThenWithResult
The promise function Then/Thenf is essentially a Success handler that waits for a subsequent promise delivery. If the original promise is not successful, the code associated with the 2nd promise is not executed.

Use **Then** when you already have a promise, or **Thenf** when you have a function that returns a promise (and takes no params).

```go
download(uri).Thenf(func() { return updateMetrics(uri)})

-- or --

// p1 and p2 are promises obtained previously
download(uri).Then(p1,p2)
```

Use **ThenWithResult** when you want to chain the result of an initial promise to the next promise:

```go
func cacheFile(result interface{}) promises.Promise {
  p := NewPromise()

  // do something with result. synch or asynch

  return p
}

{
  download(uri).ThenWithResult(cacheFile).Success(func(result interface{}) {
		file := result.([]byte)
		fmt.Printf("Downloaded %d bytes\n", len(file))
	}).Catch(func(err error) {
		fmt.Println("Error downloading/caching file: ", err)
	})
}
```

It is important to note that any use of **Then** related functions involves an intermediate promise that bridges between the intital promise and subsequent promises. In the example, the Success/Catch handlers are bound to this intermediate promise, and not directly to the promise returned from _download_, or _cacheFile_. The intermediate promise will always represent the success of _cacheFile_, but could represent the failure of either the promise from _download_ or _cacheFile_. The primary reason this matters is that had we placed a Success handler between _download_ and **ThenWithResult** is would always represent the success of the _download_ promise as the intermediate promise is not created until **ThenWithResult** is called.
