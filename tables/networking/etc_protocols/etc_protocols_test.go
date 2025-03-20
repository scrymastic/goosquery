package etc_protocols

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenEtcProtocols(t *testing.T) {
	protocols, err := GenEtcProtocols(context.Context{})
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
}
