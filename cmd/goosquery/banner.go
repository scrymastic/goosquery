package main

import "fmt"

// displayBanner shows a nice ASCII art banner
func displayBanner() {
	// ANSI color codes
	greenColor := "\033[32m"
	resetColor := "\033[0m"

	banner := `
   ____ _____  ____  _________ ___  _____  _______  __
  / __ '/ __ \/ __ \/ ___/ __ '/ / / / _ \/ ___/ / / /
 / /_/ / /_/ / /_/ (__  ) /_/ / /_/ /  __/ /  / /_/ / 
 \__, /\____/\____/____/\__, /\__,_/\___/_/   \__, /  
/____/                    /_/                /____/   -- scrymastic --                 
`
	fmt.Print(greenColor)
	fmt.Println(banner)
	fmt.Printf("GoOSQuery v%s - Windows System Information Collector\n", "0.0.1")
	fmt.Printf("Build Time: %s\n%s\n", "2025-03-18 12:00:00", resetColor)
}
