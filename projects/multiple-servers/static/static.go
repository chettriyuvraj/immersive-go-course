package static

import (
	"log"
	"net/http"
)

func Run() error {
	fs := http.Dir("../assets")

	log.Fatal(http.ListenAndServe(":8080", http.FileServer(fs)))

	return nil
}
