package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenProcesses(t *testing.T) {
	processes, err := GenProcesses()
	if err != nil {
		t.Fatalf("Failed to get processes: %v", err)
	}

	// Format and print JSON output
	jsonData, err := json.MarshalIndent(processes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal processes: %v", err)
	}
	t.Logf("Processes JSON:\n%s", string(jsonData))

	// Print number of records
	t.Logf("Number of processes: %d", len(processes))

	// Basic validation
	if len(processes) == 0 {
		t.Error("No processes returned")
	}
}

func ExampleGenProcesses() {
	processes, err := GenProcesses()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total processes: %d\n", len(processes))

	jsonData, err := json.MarshalIndent(processes, "", "  ")
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(jsonData))
}
