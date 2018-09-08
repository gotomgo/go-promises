package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	promises "github.com/gotomgo/go-promises"
)

func download(uri string) promises.Promise {
	// NewPromise returns Controller, which is a superset of Promise
	p := promises.NewPromise()

	// do the download asynchronously
	go func(p promises.Controller, uri string) {
		client := &http.Client{}
		r, _ := http.NewRequest("GET", uri, nil)

		resp, err := client.Do(r)
		if err != nil {
			// fail the promise with the error from the http client
			p.Fail(err)
			return
		}

		// keep it simple, and only succedd on 200
		if resp.StatusCode != http.StatusOK {
			p.Fail(fmt.Errorf("HTTP STATUS (%d): %s", resp.StatusCode, resp.Status))
			return
		}

		defer resp.Body.Close()

		if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			p.SucceedWithResult(bodyBytes)
		} else {
			p.Fail(err)
		}
	}(p, uri)

	return p
}

func main() {
	// because this is an example, we need something to keep the main thread alive
	var wg sync.WaitGroup

	// keep-alive until the download attempt completes
	wg.Add(1)

	download("https://github.com/gotomgo/go-promises/examples/testdata/image1.jpg").
		Success(func(result interface{}) {
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
