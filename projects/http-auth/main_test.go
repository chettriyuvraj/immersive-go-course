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
	w := httptest.NewRecorder()

	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/200", DUMMYURL), nil)
	handle200(w, getReq)
	resp := w.Result()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, resp.Body)
	require.NoError(t, err)
	require.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
	require.Equal(t, "200", buf.String())

}
