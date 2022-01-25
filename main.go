package main

import (
	"github.com/o-kos/nx-ex/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
