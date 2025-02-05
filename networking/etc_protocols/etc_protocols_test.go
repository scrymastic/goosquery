package etc_protocols

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenEtcProtocols(t *testing.T) {
	protocols, err := GenEtcProtocols()
	if err != nil {
		t.Fatalf("Failed to get protocols: %v", err)
	}

	// Verify we got some protocols
	if len(protocols) == 0 {
		t.Fatal("No protocols were returned")
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(protocols, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal protocols to JSON: %v", err)
	}
	fmt.Printf("Protocols Results:\n%s\n", string(jsonData))
	fmt.Printf("Total protocols: %d\n", len(protocols))

	// Verify some well-known protocols exist
	wellKnownProtocols := map[string]uint32{
		"tcp":  6,
		"udp":  17,
		"icmp": 1,
	}

	for name, number := range wellKnownProtocols {
		found := false
		for _, protocol := range protocols {
			if protocol.Name == name && protocol.Number == number {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Well-known protocol %s (number %d) not found in results", name, number)
		}
	}

	// Verify protocol structure
	for i, protocol := range protocols {
		if protocol.Name == "" {
			t.Errorf("Protocol at index %d has empty name", i)
		}
	}
}
