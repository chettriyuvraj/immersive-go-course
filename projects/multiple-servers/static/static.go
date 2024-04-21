package static

import (
	"fmt"
	"net/http"
)

func Run() error {
	fs := http.Dir("../assets")

	err := http.ListenAndServe(":8081", http.FileServer(fs))
	if err != nil {
		return fmt.Errorf("unable to serve file server: [%w]", err)
	}

	return nil
}
