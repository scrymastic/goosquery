package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/scrymastic/goosquery/sql/engine"
	"github.com/scrymastic/goosquery/sql/result"
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
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// Remove newline and trim whitespace
		input = strings.TrimSpace(input)

		// Handle commands
		if strings.HasPrefix(input, ".") {
			// Handle special commands
			switch strings.ToLower(input) {
			case ".exit", ".quit":
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

		// Check for exit command (legacy support)
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Exiting...")
			break
		}

		// Skip empty input
		if input == "" {
			continue
		}

		// Execute the query
		executeQuery(sqlEngine, input, jsonOutput)
	}
}

// printOutputMode prints the current output mode
func printOutputMode(jsonOutput bool) {
	if jsonOutput {
		fmt.Println("Output mode: JSON")
	} else {
		fmt.Println("Output mode: TABLE")
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
		jsonData, err := json.MarshalIndent(queryResult, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting result: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	} else {
		// Display the result as a formatted table
		displayAsTable(queryResult)
	}
}

// displayAsTable formats the query result as an ASCII table
func displayAsTable(queryResult *result.QueryResult) {
	// Unpack the query result which is a slice of map[string]interface{}
	records := *queryResult

	if len(records) == 0 {
		fmt.Println("No results found.")
		return
	}

	// Get all column names
	var columnNames []string
	columnWidths := make(map[string]int)

	// First, extract column names from the first record
	firstRecord := records[0]
	for colName := range firstRecord {
		columnNames = append(columnNames, colName)
		// Initialize column width to column name length
		columnWidths[colName] = len(colName)
	}

	// Find the maximum width needed for each column
	for _, record := range records {
		for colName, value := range record {
			valueStr := fmt.Sprintf("%v", value)
			if len(valueStr) > columnWidths[colName] {
				columnWidths[colName] = len(valueStr)
			}
		}
	}

	// Print header
	printTableDivider(columnNames, columnWidths)
	for _, colName := range columnNames {
		fmt.Printf("| %-*s ", columnWidths[colName], colName)
	}
	fmt.Println("|")
	printTableDivider(columnNames, columnWidths)

	// Print rows
	for _, record := range records {
		for _, colName := range columnNames {
			value := record[colName]
			valueStr := fmt.Sprintf("%v", value)
			fmt.Printf("| %-*s ", columnWidths[colName], valueStr)
		}
		fmt.Println("|")
	}

	// Print footer
	printTableDivider(columnNames, columnWidths)
	fmt.Printf("%d rows in set\n", len(records))
}

// printTableDivider prints a horizontal line for the table
func printTableDivider(columnNames []string, columnWidths map[string]int) {
	for _, colName := range columnNames {
		fmt.Printf("+-%s-", strings.Repeat("-", columnWidths[colName]))
	}
	fmt.Println("+")
}
