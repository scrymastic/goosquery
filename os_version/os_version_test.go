package os_version

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetOSVersion(t *testing.T) {
	osVersion, err := GenOSVersion()
	if err != nil {
		t.Fatalf("Failed to get OS version: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(osVersion, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal OS version to JSON: %v", err)
	}
	fmt.Printf("OS Version Information:\n%s\n", string(jsonData))
}
