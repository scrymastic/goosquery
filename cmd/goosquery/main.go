package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/scrymastic/goosquery/sql/engine"
)

func main() {
	// Parse command-line flags
	queryFlag := flag.String("q", "", "SQL query to execute")
	interactiveFlag := flag.Bool("i", false, "Run in interactive mode")
	jsonFlag := flag.Bool("json", false, "Output results in JSON format")
	flag.Parse()

	// Create SQL engine
	sqlEngine := engine.NewEngine()

	// If interactive mode specified or no query provided, start interactive mode
	if *interactiveFlag || *queryFlag == "" {
		runInteractiveMode(sqlEngine, *jsonFlag)
		return
	}

	// Execute a single query in non-interactive mode
	executeQuery(sqlEngine, *queryFlag, *jsonFlag)
}

// runInteractiveMode starts an interactive REPL for executing SQL queries
func runInteractiveMode(sqlEngine *engine.Engine, jsonOutput bool) {
	displayBanner()

	reader := bufio.NewReader(os.Stdin)

	// Show current output mode
	printOutputMode(jsonOutput)

	for {
		fmt.Print("goosquery> ")

		// Read user input (command or SQL)
		input, isCommand, err := readSQLInput(reader)
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// Skip if input is empty
		if input == "" {
			continue
		}

		// Handle commands
		if isCommand {
			// Handle special commands
			switch strings.ToLower(input) {
			case ".quit":
				fmt.Println("Exiting...")
				return
			case ".json":
				jsonOutput = true
				fmt.Println("Output mode set to JSON")
				continue
			case ".table":
				jsonOutput = false
				fmt.Println("Output mode set to TABLE")
				continue
			case ".mode":
				printOutputMode(jsonOutput)
				continue
			case ".help":
				fmt.Println("Commands:")
				fmt.Println("  .quit        - Exit the program")
				fmt.Println("  .json        - Switch to JSON output mode")
				fmt.Println("  .table       - Switch to table output mode")
				fmt.Println("  .mode        - Show current output mode")
				fmt.Println("  .help        - Show this help message")
				continue
			default:
				fmt.Printf("Unknown command: %s\n", input)
				continue
			}
		}

		// Execute the query
		executeQuery(sqlEngine, input, jsonOutput)
	}
}

// executeQuery runs a SQL query and displays the result
func executeQuery(sqlEngine *engine.Engine, query string, jsonOutput bool) {
	queryResult, err := sqlEngine.Execute(query)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	if jsonOutput {
		// Convert result to JSON and print
		jsonData, err := formatResultAsJSON(queryResult)
		if err != nil {
			fmt.Printf("Error formatting result: %v\n", err)
			return
		}
		fmt.Println(jsonData)
		fmt.Printf("Total rows: %d\n", queryResult.Size())
	} else {
		// Display the result as a formatted table
		displayAsTable(queryResult)
		fmt.Printf("Total rows: %d\n", queryResult.Size())
	}
}
