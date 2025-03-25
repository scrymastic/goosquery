package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
)

// readSQLInput reads input from the user, either a command or a multi-line SQL query
// Returns the input string, whether it's a command, and any error encountered
func readSQLInput(reader *bufio.Reader) (string, bool, error) {
	var sqlBuffer strings.Builder

	// Read first line
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}

	// Trim whitespace
	line = strings.TrimSpace(line)

	// Check if this is a special command
	if strings.HasPrefix(line, ".") {
		return line, true, nil
	}

	// Not a command, start collecting SQL
	sqlBuffer.WriteString(line)

	// If line ends with semicolon, we're done
	if strings.HasSuffix(line, ";") {
		return sqlBuffer.String(), false, nil
	}

	// Continue reading until semicolon is found
	for {
		if len(sqlBuffer.String()) > 0 {
			fmt.Print(".........> ")
		} else {
			fmt.Print("goosquery> ")
		}
		line, err = reader.ReadString('\n')
		if err != nil {
			return sqlBuffer.String(), false, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		sqlBuffer.WriteString(" " + line)

		// Check if we've reached the end of the SQL statement
		if strings.HasSuffix(line, ";") {
			break
		}
	}

	return sqlBuffer.String(), false, nil
}

// printOutputMode prints the current output mode
func printOutputMode(jsonOutput bool) {
	if jsonOutput {
		fmt.Println("Output mode: JSON")
	} else {
		fmt.Println("Output mode: TABLE")
	}
}

// displayAsTable formats the query result as an ASCII table
func displayAsTable(queryResult *result.Results) {
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

// formatResultAsJSON converts the query result to JSON format
func formatResultAsJSON(queryResult *result.Results) (string, error) {
	jsonData, err := json.MarshalIndent(queryResult, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
