package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	URL           = "http://localhost:8080"
	RETRYAFTERKEY = "Retry-After"
)

var ClientUnexpectedError = errors.New("unexpected error making client request")

func main() {
	for {
		fmt.Printf("\n\nMaking new request...")
		respStr, err := MakeRequest()
		if err != nil {
			if errors.Is(err, ClientUnexpectedError) {
				fmt.Fprintf(os.Stderr, "\nerror: irrevcoverable error: %v", err)
				os.Exit(1)
				continue
			}
			fmt.Fprintf(os.Stderr, "\nerror: %v", err)
			continue
		}
		fmt.Printf("\n%v", respStr)
		time.Sleep(3 * time.Second)
	}
}

func MakeRequest() (string, error) {
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(URL)
	if err != nil { /* request itself returns an error */
		return "", ClientUnexpectedError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests { /* Non-2xx response codes wont error out, but this is 429 */
		return handleRequestRetry(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

/**** HELPERS ****/
func handleRequestRetry(resp *http.Response) (string, error) {
	retryVal := resp.Header.Get(RETRYAFTERKEY)
	secsUntilRetry, err := strconv.Atoi(retryVal) /* Sometimes the server sends seconds to wait */
	if err != nil {
		retryTime, err := time.Parse(http.TimeFormat, retryVal) /* Sometimes the server sends exact retry time */
		if err != nil {
			return "", fmt.Errorf("invalid format for retry time")
		}

		timeUntilRetry := time.Until(retryTime)
		if timeUntilRetry > 0 { /* If retry time has not already passed by the time response arrives  */
			fmt.Printf("\nSleeping for %v nanoseconds...", timeUntilRetry)
			time.Sleep(timeUntilRetry)
		}
		fmt.Printf("\nRetrying...")
		return MakeRequest()
	}

	fmt.Printf("\nSleeping for %v seconds...", secsUntilRetry)
	time.Sleep(time.Second * time.Duration(secsUntilRetry))
	fmt.Printf("\nRetrying...")
	return MakeRequest()
}
