package fetcher

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
	RETRYAFTERKEY = "Retry-After"
	TIMEOUT       = 5
)

/**** ERRORS ****/

var ClientUnexpectedError = errors.New("unexpected error making client request")

type RetryError struct {
	msg     string
	retryIn time.Duration
}

func (r *RetryError) Error() string {
	return r.msg
}

/**** FETCHER ****/

type Fetcher struct {
	client http.Client
}

func NewFetcher(rt http.RoundTripper) *Fetcher {
	return &Fetcher{
		client: http.Client{
			Transport: rt,
			Timeout:   TIMEOUT * time.Second,
		},
	}
}

func (f *Fetcher) MakeRequest(URL string) (string, error) {
	fmt.Fprintf(os.Stderr, "\n\nMaking new request to %v...", URL)

	resp, err := f.client.Get(URL)
	if err != nil { /* request itself returns an error */
		return "", ClientUnexpectedError
	}
	defer resp.Body.Close()

	switch sc := resp.StatusCode; sc {
	case http.StatusOK:
		return f.handle200(resp)
	case http.StatusTooManyRequests:
		return f.handle429(resp)
	default:
		return f.handleUnexpectedResponse(resp)
	}
}

/**** HELPERS ****/

func (f *Fetcher) handle200(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func (f *Fetcher) handle429(resp *http.Response) (string, error) {
	retryVal := resp.Header.Get(RETRYAFTERKEY)
	timeUntilRetry, err := f.parseDelay(retryVal)
	if err != nil {
		return "", fmt.Errorf("error parsing retry-time delay: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\nSleeping for %v", timeUntilRetry)
	time.Sleep(timeUntilRetry)
	return "", &RetryError{msg: "retry the request", retryIn: timeUntilRetry}
}

func (f *Fetcher) handleUnexpectedResponse(resp *http.Response) (string, error) {
	/* Different from UnexpectedError, this is an actual response that is unexpected */
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading unexpected response body")
	}
	return fmt.Sprintf("status code: %d, response body: %v", resp.StatusCode, string(body)), nil
}

func (f *Fetcher) parseDelay(retryVal string) (time.Duration, error) {
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
