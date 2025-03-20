package etc_hosts

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenEtcHosts(t *testing.T) {
	entries, err := GenEtcHosts(context.Context{})
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
		if address, ok := entry["address"]; !ok || address.(string) == "" {
			t.Errorf("Entry %d has empty or missing address", i)
		}
		if hostnames, ok := entry["hostnames"]; !ok || hostnames.(string) == "" {
			t.Errorf("Entry %d has empty or missing hostnames", i)
		}
	}
}
