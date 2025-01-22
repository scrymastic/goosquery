package main

import (
	"fmt"
	"runtime"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Hello, World!")
	}
}
