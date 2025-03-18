package logical_drives

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenLogicalDrives(t *testing.T) {
	drives, err := GenLogicalDrives()
	if err != nil {
		t.Fatalf("Failed to get logical drives: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(drives, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal logical drives to JSON: %v", err)
	}
	fmt.Printf("Logical Drives Results:\n%s\n", string(jsonData))
	fmt.Printf("Total drives: %d\n", len(drives))
}
