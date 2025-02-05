package etc_hosts

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenEtcHosts(t *testing.T) {
	entries, err := GenEtcHosts()
	if err != nil {
		t.Fatalf("Failed to get hosts entries: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal hosts entries to JSON: %v", err)
	}
	fmt.Printf("Hosts File Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", len(entries))

	// Basic validation of entries
	for i, entry := range entries {
		if entry.Address == "" {
			t.Errorf("Entry %d has empty address", i)
		}
		if entry.Hostnames == "" {
			t.Errorf("Entry %d has empty hostnames", i)
		}
	}
}
