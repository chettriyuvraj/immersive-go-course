package fetcher

import (
	"net/http"
	"testing"
	"time"

	"github.com/chettriyuvraj/immersive-go-course/projects/output-and-error-handling/client/testutils"
	"github.com/stretchr/testify/require"
)

func TestParseDelay(t *testing.T) {
	mockTime := time.Date(2035, time.March, 5, 14, 10, 2, 0, time.UTC)
	mockTimePast := time.Date(2021, time.March, 5, 14, 10, 2, 0, time.UTC)
	mockTimeString := "Mon, 05 Mar 2035 14:10:02 GMT"
	mockTimeStringPast := "Mon, 05 Mar 2021 14:10:02 GMT"
	type args struct {
		retryVal string
	}
	tests := []struct {
		name    string
		args    args
		delay   time.Duration
		wantErr bool
		errStr  string
	}{
		{
			name: "valid delay seconds",
			args: args{
				retryVal: "5",
			},
			delay: 5 * time.Second,
		},
		{
			name: "valid delay timestamp",
			args: args{
				retryVal: mockTimeString,
			},
			delay: time.Until(mockTime),
		},
		{
			name: "valid delay timestep - past",
			args: args{
				retryVal: mockTimeStringPast,
			},
			delay: time.Until(mockTimePast),
		},
		{
			name: "invalid delay",
			args: args{
				retryVal: "bull",
			},
			wantErr: true,
			errStr:  "invalid format for retry time: bull",
		},
		{
			name: "invalid delay timestep",
			args: args{
				retryVal: "Wednesday, 32 February 2051 14:00:01 GM",
			},
			wantErr: true,
			errStr:  "invalid format for retry time: Wednesday, 32 February 2051 14:00:01 GM",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := NewFetcher(http.DefaultTransport)
			got, err := f.parseDelay(tc.args.retryVal)

			/* Let's make sure err exists only when there is a wantErr in our testCase - no mismatch */
			require.Equal(t, (err != nil), tc.wantErr)

			/* Else isn't idiomatic go apparently but it suits test cases like these pretty well */
			if tc.wantErr {
				require.EqualError(t, err, tc.errStr)
			} else {
				require.InDelta(t, got/time.Second, tc.delay/time.Second, 1)
			}
		})
	}
}

func TestHandle200(t *testing.T) {
	body := "Today it will be sunny!"
	mockRoundTripper := testutils.NewMockRoundTripper(t)
	f := NewFetcher(mockRoundTripper)
	mockRoundTripper.StubResponse(200, &http.Header{}, body)
	resp, err := f.MakeRequest("http://www.dummyurl.com")
	require.NoError(t, err)
	require.NotEqual(t, resp, body)
}
