package processes

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

	fmt.Printf("Processes Results:\n%s\n", string(jsonData))
	fmt.Printf("Total processes: %d\n", len(processes))
}
