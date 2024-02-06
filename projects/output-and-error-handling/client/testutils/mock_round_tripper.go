package testutils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockRoundTripper struct {
	t         *testing.T
	responses []*http.Response
	reqIndex  int
}

func NewMockRoundTripper(t *testing.T) *MockRoundTripper {
	return &MockRoundTripper{
		t: t,
	}
}

func (m *MockRoundTripper) StubResponse(statusCode int, header *http.Header, body string) {
	resp := &http.Response{
		StatusCode: statusCode,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     *header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	m.responses = append(m.responses, resp)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.t.Helper()

	require.Less(m.t, m.reqIndex, len(m.responses), fmt.Sprintf("error: number of requests %d, number of responses %d", m.reqIndex+1, len(m.responses)))
	return m.responses[m.reqIndex], nil
}
