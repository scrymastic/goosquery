package interface_addresses

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenInterfaceAddresses(t *testing.T) {
	addresses, err := GenInterfaceAddresses(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get interface addresses: %v", err)
	}

	// Verify we got at least one interface
	if addresses.Size() == 0 {
		t.Error("Expected at least one interface address, got none")
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(addresses, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal interface addresses to JSON: %v", err)
	}
	fmt.Printf("Interface Addresses Results:\n%s\n", string(jsonData))
	fmt.Printf("Total interfaces: %d\n", addresses.Size())
}
