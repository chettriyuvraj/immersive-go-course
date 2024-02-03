package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chettriyuvraj/immersive-go-course/projects/output-and-error-handling/client/fetcher"
)

const URL = "http://localhost:8080"

func main() {
	f := fetcher.NewFetcher(http.DefaultTransport)
	for {
		respStr, err := f.MakeRequest(URL)
		if err != nil {

			switch {
			case errors.As(err, &fetcher.ClientUnexpectedError):
				fmt.Fprintf(os.Stderr, "\nerror: irrevcoverable error: %v", err)
				os.Exit(1)
			case errors.As(err, &fetcher.RetryError{}):
				fmt.Fprintf(os.Stderr, "\nerror: retrying", err)
				continue
			}
		}
		fmt.Printf("\n%v", respStr)
		time.Sleep(2 * time.Second)
	}
}
