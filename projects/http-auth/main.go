package main

import "net/http"

func main() {

	http.HandleFunc("/", http.NotFoundHandler().ServeHTTP)

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
	})

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	http.ListenAndServe(":8080", nil)
}
