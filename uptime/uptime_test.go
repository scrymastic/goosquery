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

	t.Logf("System Uptime JSON:\n%s", string(jsonData))
	t.Logf("Number of records: %d", len(uptime))
}

func ExampleGenUptime() {
	uptime, err := GenUptime()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total records: %d\n", len(uptime))

	jsonData, err := json.MarshalIndent(uptime, "", "  ")
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(jsonData))
}
