package uptime

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenUptime(t *testing.T) {
	uptime, err := GenUptime()
	if err != nil {
		t.Fatalf("Failed to get uptime: %v", err)
	}

	jsonData, err := json.MarshalIndent(uptime, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal uptime: %v", err)
	}

	fmt.Printf("Uptime Results:\n%s\n", string(jsonData))
}
