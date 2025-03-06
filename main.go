package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Check if the os is windows amd64
	if runtime.GOOS != "windows" || runtime.GOARCH != "amd64" {
		fmt.Println("This program only supports Windows AMD64")
		os.Exit(1)
	}

	fmt.Println("Running on Windows AMD64")
}
