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

func TestInvalidURLs(t *testing.T) {

	type args struct {
		urlPath string
	}
	type resp struct {
		respStr    string
		statusCode int
	}
	var invalidURLResp resp = resp{respStr: "404 page not found\n", statusCode: http.StatusNotFound}

	tests := []struct {
		args args
		resp resp
	}{
		{
			args: args{
				urlPath: "/554",
			},
			resp: invalidURLResp,
		},
		{
			args: args{
				urlPath: "/ddsfg",
			},
			resp: invalidURLResp,
		},
	}
	for _, tc := range tests {
		t.Run("test invalid path", func(t *testing.T) {
			url := fmt.Sprintf("%s%s", DUMMYURL, tc.args.urlPath)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()
			handleBase(w, req)
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

func TestHandleAuthenticated(t *testing.T) {
	url := fmt.Sprintf("%s/authenticated", DUMMYURL)
	type args struct {
		req     *http.Request
		headers map[string]string
	}
	type resp struct {
		respStr    string
		statusCode int
		headers    map[string]string
	}
	tests := []struct {
		name string
		args args
		resp resp
	}{
		{
			name: "request without credentials",
			args: args{
				req: httptest.NewRequest(http.MethodGet, url, nil),
			},
			resp: resp{
				statusCode: http.StatusUnauthorized,
				headers: map[string]string{
					"WWW-Authenticate": `Basic realm = "Access to the staging site"`,
				},
			},
		},
		{
			name: "request with valid credentials",
			args: args{
				req: httptest.NewRequest(http.MethodGet, url, nil),
				headers: map[string]string{
					"Authorization": "Basic Y2hldHRyaXl1dnJhajpraW5n", //"chettriyuvraj:king"
				},
			},
			resp: resp{
				respStr:    "<!DOCTYPE html>\n<html>Hello chettriyuvraj!",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "request with invalid credentials",
			args: args{
				req: httptest.NewRequest(http.MethodGet, url, nil),
				headers: map[string]string{
					"Authorization": "Basic G2hldHRyaXl1dnJhajpraW5n", //"chettriyuvraj:king"
				},
			},
			resp: resp{
				respStr:    "Credentials wrong or invalid",
				statusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			for k, v := range tc.args.headers {
				tc.args.req.Header.Set(k, v)
			}
			handleAuthenticated(w, tc.args.req)

			resp := w.Result()
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, resp.Body)
			require.NoError(t, err)
			require.Equal(t, tc.resp.respStr, buf.String())
			require.Equal(t, tc.resp.statusCode, resp.StatusCode)
			require.Equal(t, len(tc.resp.headers), len(resp.Header))
			for k := range tc.resp.headers {
				require.Equal(t, tc.resp.headers[k], resp.Header.Get(k))
			}
		})
	}
}
