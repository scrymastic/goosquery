package dns_cache

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenDNSCache(t *testing.T) {
	cache, err := GenDNSCache()
	if err != nil {
		t.Fatalf("Failed to get DNS cache: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal DNS cache: %v", err)
	}
	fmt.Printf("DNS Cache Results:\n%s\n", string(jsonData))
	fmt.Printf("Total DNS cache entries: %d\n", len(cache))
}

func ExampleGenDNSCache() {
	cache, err := GenDNSCache()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Number of DNS cache entries: %d\n", len(cache))
}
