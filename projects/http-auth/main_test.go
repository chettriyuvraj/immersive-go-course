package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const DUMMYURL = "http://localhost:8080"

func TestHandle200(t *testing.T) {
	url := fmt.Sprintf("%s/200", DUMMYURL)
	type args struct {
		req *http.Request
	}
	type resp struct {
		respStr    string
		statusCode int
	}

	tests := []struct {
		name string
		args args
		resp resp
	}{
		{
			name: "valid response",
			args: args{
				req: httptest.NewRequest(http.MethodGet, url, nil),
			},
			resp: resp{
				respStr:    "200",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "404 response",
			args: args{
				req: httptest.NewRequest(http.MethodPost, url, nil),
			},
			resp: resp{
				respStr:    "404 page not found\n",
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handle200(w, tc.args.req)
			resp := w.Result()
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)
			require.Equal(t, tc.resp.respStr, buf.String())
			require.Equal(t, tc.resp.statusCode, resp.StatusCode)
		})
	}
}
