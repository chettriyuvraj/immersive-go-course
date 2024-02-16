package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleBase)
	http.HandleFunc("/authenticated", handleAuthenticated)
	http.HandleFunc("/200", handle200)
	http.HandleFunc("/500", handle500)
	http.ListenAndServe(":8080", nil)
}

func handleBase(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte("<em>Hello world</em>\n"))
	buf.Write([]byte("<ul>\n"))
	queryParams := r.URL.Query()
	for k, vList := range queryParams {
		v := vList[0]
		v = html.EscapeString(v)
		buf.Write([]byte(fmt.Sprintf("<li>%s: %s</li>\n", k, v)))
	}

	w.Header().Set("Content-Type", "text/html")
	buf.Write([]byte("</ul>\n"))
	_, err := w.Write(buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func handleAuthenticated(w http.ResponseWriter, r *http.Request) {
	username, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Basic realm = "Access to the staging site"`)
		return
	}
	if username == "chettriyuvraj" && pass == "king" {
		w.Write([]byte("<!DOCTYPE html>\n<html>Hello username!"))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Credentials wrong or invalid"))
}

func handle200(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("200"))
}

func handle500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}
