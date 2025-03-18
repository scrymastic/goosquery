package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/scrymastic/goosquery/sql/engine"
)

func main() {
	// Parse command-line flags
	queryFlag := flag.String("q", "", "SQL query to execute")
	flag.Parse()

	// Set default query if none provided
	query := *queryFlag
	if query == "" {
		query = "SELECT * FROM processes WHERE pid = 4"
	}

	// Create SQL engine
	sqlEngine := engine.NewEngine()

	// Execute the query
	result, err := sqlEngine.Execute(query)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	// Convert result to JSON and print
	jsonOutput, err := result.ToJSON()
	if err != nil {
		log.Fatalf("Error formatting result: %v", err)
	}
	fmt.Println(jsonOutput)
}
