package programs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenPrograms(t *testing.T) {
	programs, err := GenPrograms()
	if err != nil {
		t.Fatalf("Failed to get programs: %v", err)
	}

	// Format and print JSON output
	jsonData, err := json.MarshalIndent(programs, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal programs: %v", err)
	}
	t.Logf("Installed Programs JSON:\n%s", string(jsonData))

	// Print number of records
	t.Logf("Number of programs: %d", len(programs))

	// Basic validation
	if len(programs) == 0 {
		t.Error("No programs returned")
	}
}

func ExampleGenPrograms() {
	programs, err := GenPrograms()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total installed programs: %d\n", len(programs))

	jsonData, err := json.MarshalIndent(programs, "", "  ")
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(jsonData))
}
