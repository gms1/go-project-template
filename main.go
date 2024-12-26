package main

import (
	"os"

	"github.com/gms1/go-project-template/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
