package main

import (
	"fmt"
	"os"
	"servers/api"
)

func main() {
	err := api.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in serving static content: [%w]", err)
	}
}
