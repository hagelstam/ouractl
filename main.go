package main

import (
	"os"

	"github.com/hagelstam/ouractl/cmd/ouractl"
)

func main() {
	if err := ouractl.Execute(); err != nil {
		os.Exit(1)
	}
}
