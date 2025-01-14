package authenticode

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenAuthenticode(t *testing.T) {
	// Test with Windows system file
	testPath := `C:\Windows\System32\ntoskrnl.exe`
	entries, err := GenAuthenticode(testPath)
	if err != nil {
		t.Fatalf("Failed to get authenticode info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal authenticode info to JSON: %v", err)
	}
	fmt.Printf("Authenticode Results for %s:\n%s\n", testPath, string(jsonData))
	fmt.Printf("Total entries: %d\n", len(entries))
}
