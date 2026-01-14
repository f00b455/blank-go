package main

import (
	"fmt"
	"os"

	"github.com/f00b455/blank-go/cmd/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
