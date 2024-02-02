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

	switch sc := resp.StatusCode; sc {
	case http.StatusOK:
		return handle200(resp)
	case http.StatusTooManyRequests:
		return handle429(resp)
	default:
		return handleUnexpectedResponse(resp)
	}
}

/**** HELPERS ****/

func handle200(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func handle429(resp *http.Response) (string, error) {
	retryVal := resp.Header.Get(RETRYAFTERKEY)
	timeUntilRetry, err := parseDelay(retryVal)
	if err != nil {
		return "", fmt.Errorf("error parsing retry-time delay: %w", err)
	}

	fmt.Printf("\nSleeping for %v", timeUntilRetry)
	time.Sleep(timeUntilRetry)
	return MakeRequest()
}

func handleUnexpectedResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading unexpected response body")
	}
	return fmt.Sprintf("status code: %d, response body: %v", resp.StatusCode, string(body)), nil
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
