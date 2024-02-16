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
	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func handleBase(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	username, pass, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm = "Access to the staging site"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if username == "chettriyuvraj" && pass == "king" {
		w.Write([]byte("<!DOCTYPE html>\n<html>Hello chettriyuvraj!"))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Write([]byte("Credentials wrong or invalid"))
	w.WriteHeader(http.StatusUnauthorized)
}

func handle200(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("200"))
	w.WriteHeader(http.StatusOK)
}

func handle500(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Internal Server Error"))
	w.WriteHeader(http.StatusInternalServerError)
}
