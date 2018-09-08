# GO Promises package

go-promises is a flexibile, easy to use implementation of promises written in GO. This version is heavily based on a C# implementation I wrote for use with the Unity 3D game engine, which has been used / tested extensively in successful commercial products. The primary difference in the implementations (beyond language) is that Unity 3D is primarily single-threaded with a co-routine model, whereas the GO implementation must support concurrent operations. For those interested in the C# version, it can be found here: https://bitbucket.org/codemules/promises

## Getting Started

1. Download and install it:

```sh
$ go get github.com/gotomgo/go-promises
```

2. Import it in your code:

```go
import promises "github.com/gotomgo/promises"
```
