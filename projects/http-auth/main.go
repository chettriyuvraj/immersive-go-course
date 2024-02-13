package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
			return
		}

		b := make([]byte, 2048)
		buf := new(bytes.Buffer)
		body := r.Body
		for {
			n, err := body.Read(b)
			fmt.Println(string(b))
			if err != nil && err != io.EOF {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err2 := buf.Write(b[:n])
			if err2 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err == io.EOF {
				break
			}
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
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
