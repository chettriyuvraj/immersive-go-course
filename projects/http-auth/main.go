package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("200"))
	})

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	http.ListenAndServe(":8080", nil)
}
