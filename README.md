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

  // perform work here. Synchronously or Asynchronously

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
    if p.Succeeded() {
      result := p.Result()
    } else {
      fmt.Println("Error from myFunction: ",p.Error())
    }
  })
```
