package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	promises "github.com/gotomgo/go-promises"
)

// doDownload does the actual download, synch or asynch
//
//  Notes
//    The function requires a promises.Controller as it will deliver the
//    promise, whereas a normal consumer of a promise only needs
//    and instance of promises.Promise
//
//    The function returns the Promise *only* to make the code
//    a little cleaner (don't need seperate line with 'return'.)
//    We could re-structure the code to avoid this, but its an example
func doDownload(uri string, p promises.Controller) promises.Promise {
	client := &http.Client{}
	r, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		// fail the promise with the error from the http request
		return p.Fail(err)
	}

	resp, err := client.Do(r)
	if err != nil {
		// fail the promise with the error from the http client
		return p.Fail(err)
	}

	// keep it simple, and only succedd on 200
	if resp.StatusCode != http.StatusOK {
		return p.Fail(fmt.Errorf("HTTP STATUS (%d): %s", resp.StatusCode, resp.Status))
	}

	defer resp.Body.Close()

	// JIC U R Curious
	// fmt.Println("Download content length = ", resp.ContentLength)

	if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
		return p.SucceedWithResult(bodyBytes)
	} else {
		return p.Fail(err)
	}
}

// synchDownload downloads a resource and returns a promise
// for delivery of the results. Because the download is synchronous,
// the promise is delivered prior to the function returning
func synchDownload(uri string) promises.Promise {
	// NewPromise returns Controller, which is a superset of Promise
	p := promises.NewPromise()

	doDownload(uri, p)

	return p
}

// aynchDownload starts a GO routine to download the file and returns a promise
// for delivery of the results
func asynchDownload(uri string) promises.Promise {
	// NewPromise returns Controller, which is a superset of Promise
	p := promises.NewPromise()

	// do the download asynchronously
	go func(uri string, p promises.Controller) {
		doDownload(uri, p)
	}(uri, p)

	return p
}

func main() {
	// because this is an example, we need something to keep the main thread alive
	var wg sync.WaitGroup

	// keep-alive until the download attempt completes
	wg.Add(1)

	uri := "https://github.com/gotomgo/go-promises/tree/master/examples/testdata/image1.jpg"

	// use a function pointer so we can switch between asynch (default) or
	// synch download
	var download func(string) promises.Promise

	// default is asynch
	download = asynchDownload

	// allow a command-line arg to tell us to do the download synchronously
	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "synch" {
		download = synchDownload
	}

	// do the download and process the results
	download(uri).Success(func(result interface{}) {
		bodyBytes := result.([]byte)
		fmt.Printf("Downloaded %d bytes\n", len(bodyBytes))
	}).Catch(func(err error) {
		fmt.Println("Error downloading file: ", err)
	}).Always(func(p promises.Controller) {
		// tell the main thread the download attempt completed
		wg.Done()
	})

	// wait for the download attempt to complete, then exit
	wg.Wait()
}
