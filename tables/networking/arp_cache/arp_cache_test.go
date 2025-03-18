package arp_cache

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenARPCache(t *testing.T) {
	entries, err := GenARPCache()
	if err != nil {
		t.Fatalf("Failed to get ARP entries: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ARP entries to JSON: %v", err)
	}
	fmt.Printf("ARP Cache Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", len(entries))
}
