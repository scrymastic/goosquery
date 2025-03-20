package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/scrymastic/goosquery/sql/engine"
	"github.com/scrymastic/goosquery/sql/result"
)

func TestQuery(t *testing.T) {
	sqlEngine := engine.NewEngine()
	query := "SELECT * FROM processes;"
	queryResult, err := sqlEngine.Execute(query)
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}

	jsonData, err := json.MarshalIndent(queryResult, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal query result: %v", err)
	}

	fmt.Println(string(jsonData))
}

func TestDisplayAsTable(t *testing.T) {
	// Create a sample query result
	queryResult := &result.QueryResult{
		{
			"name":  "test1",
			"value": int64(123),
		},
		{
			"name":  "test2",
			"value": int64(456),
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	displayAsTable(queryResult)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that the output contains expected table elements
	if !strings.Contains(output, "name") || !strings.Contains(output, "value") {
		t.Errorf("Table headers not found in output")
	}
	if !strings.Contains(output, "test1") || !strings.Contains(output, "test2") {
		t.Errorf("Table data not found in output")
	}
	if !strings.Contains(output, "123") || !strings.Contains(output, "456") {
		t.Errorf("Table values not found in output")
	}
	if !strings.Contains(output, "2 rows in set") {
		t.Errorf("Row count not found in output")
	}
}

func TestPrintTableDivider(t *testing.T) {
	columnNames := []string{"col1", "col2"}
	columnWidths := map[string]int{
		"col1": 4,
		"col2": 4,
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	printTableDivider(columnNames, columnWidths)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that the output contains expected divider characters
	if !strings.Contains(output, "+") || !strings.Contains(output, "-") {
		t.Errorf("Table divider not formatted correctly")
	}
}

func TestDisplayAsTableWithEmptyResults(t *testing.T) {
	// Create an empty query result
	queryResult := &result.QueryResult{}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	displayAsTable(queryResult)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that the output contains the empty result message
	if !strings.Contains(output, "No results found") {
		t.Errorf("Empty result message not found in output")
	}
}

func TestPrintOutputMode(t *testing.T) {
	testCases := []struct {
		jsonOutput bool
		expected   string
	}{
		{true, "Output mode: JSON"},
		{false, "Output mode: TABLE"},
	}

	for _, tc := range testCases {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call the function
		printOutputMode(tc.jsonOutput)

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, tc.expected) {
			t.Errorf("Expected output to contain %q, got %q", tc.expected, output)
		}
	}
}
