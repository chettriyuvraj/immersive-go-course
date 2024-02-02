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
		respStr, err := MakeRequest()
		if err != nil {
			if errors.Is(err, ClientUnexpectedError) {
				fmt.Fprintf(os.Stderr, "\nerror: irrevcoverable error: %v", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "\nerror: %v", err)
			continue
		}
		fmt.Printf("\n%v", respStr)
		time.Sleep(3 * time.Second)
	}
}

func MakeRequest() (string, error) {
	fmt.Printf("\n\nMaking new request...")
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
	timeUntilRetry, err := parseDelay(retryVal)
	if err != nil {
		return "", fmt.Errorf("error parsing retry-time delay: %w", err)
	}

	fmt.Printf("\nSleeping for %v seconds...", timeUntilRetry/time.Second)
	time.Sleep(timeUntilRetry)
	return MakeRequest()
}

func parseDelay(retryVal string) (time.Duration, error) {
	/* Sometimes the server sends seconds to wait  */
	secsUntilRetry, err := strconv.Atoi(retryVal)

	/* Sometimes the server sends exact retry time */
	if err != nil {
		retryTime, err := time.Parse(http.TimeFormat, retryVal)
		if err != nil {
			return time.Second, fmt.Errorf("invalid format for retry time: %v", retryVal)
		}

		return time.Until(retryTime), nil
	}

	return time.Second * time.Duration(secsUntilRetry), nil
}
