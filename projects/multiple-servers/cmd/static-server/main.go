package main

import (
	"fmt"
	"os"
	"servers/static"
)

func main() {
	err := static.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in serving static content: [%w]", err)
	}
}
