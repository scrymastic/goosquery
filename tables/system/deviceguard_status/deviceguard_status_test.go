package deviceguard_status

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenDeviceguardStatus(t *testing.T) {
	status, err := GenDeviceguardStatus()
	if err != nil {
		t.Fatalf("Failed to get Device Guard status: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Device Guard status: %v", err)
	}
	fmt.Printf("Device Guard Status:\n%s\n", string(jsonData))
	fmt.Printf("Total records: %d\n", len(status))
}
