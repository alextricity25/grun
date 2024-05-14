package main

import (
	"fmt"
	"os"
)

func showBanner() {
	fmt.Fprintf(os.Stderr, header, version)
}

func showUsage() {
	main := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", main)
	fmt.Fprint(os.Stderr, options)
	fmt.Fprintf(os.Stderr, examples, main, main)
}
