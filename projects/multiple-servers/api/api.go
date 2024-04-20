package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

var images []Image = []Image{
	{
		Title:   "Sunset",
		AltText: "Clouds at sunset",
		URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
	},
	{
		Title:   "Mountain",
		AltText: "A mountain at sunset",
		URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
	},
}

func Run() error {

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		enc := json.NewEncoder(w)
		enc.Encode(images)
	})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return fmt.Errorf("unable to serve file server: [%w]", err)
	}

	return nil
}
