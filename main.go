// Package main is the cgen application.
package main

import (
	"os"
)

var version string

func main() {
	root := NewRootCmd()
	root.Version = version

	if err := root.Execute(); err != nil {
		os.Exit(-1)
	}
}

