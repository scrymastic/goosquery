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

// displayAsTable prints the query result in a formatted table with box-drawing characters
func displayAsTable(queryResult *result.Results) {
	if queryResult == nil || queryResult.Size() == 0 {
		return
	}

	// Get column names
	columns := (*queryResult).GetColumns()

	// Calculate column widths
	columnWidths := make(map[string]int)
	for _, col := range columns {
		columnWidths[col] = len(col)
	}

	// Find the maximum width needed for each column
	for _, row := range *queryResult {
		for col, value := range row {
			strValue := fmt.Sprintf("%v", value)
			if len(strValue) > columnWidths[col] {
				columnWidths[col] = len(strValue)
			}
		}
	}

	// Print header
	// Top border
	fmt.Print("┌")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", columnWidths[col]+2))
		if i < len(columns)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐")

	// Column names
	fmt.Print("│")
	for _, col := range columns {
		fmt.Printf(" %-*s │", columnWidths[col], col)
	}
	fmt.Println()

	// Separator line
	fmt.Print("├")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", columnWidths[col]+2))
		if i < len(columns)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")

	// Print rows
	for _, row := range *queryResult {
		fmt.Print("│")
		for _, col := range columns {
			value, exists := row[col]
			strValue := ""
			if exists {
				strValue = fmt.Sprintf("%v", value)
			}
			fmt.Printf(" %-*s │", columnWidths[col], strValue)
		}
		fmt.Println()
	}

	// Bottom border
	fmt.Print("└")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", columnWidths[col]+2))
		if i < len(columns)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("┘")
}

// formatResultAsJSON converts the query result to JSON format
func formatResultAsJSON(queryResult *result.Results) (string, error) {
	jsonData, err := json.MarshalIndent(queryResult, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
