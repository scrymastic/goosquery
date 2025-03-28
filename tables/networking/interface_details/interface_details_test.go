package interface_details

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenInterfaceDetails(t *testing.T) {
	interfaces, err := GenInterfaceDetails(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get interface details: %v", err)
	}

	// Verify we got at least one interface
	if interfaces.Size() == 0 {
		t.Error("No network interfaces found")
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(interfaces, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal interface details to JSON: %v", err)
	}
	fmt.Printf("Interface Details Results:\n%s\n", string(jsonData))
	fmt.Printf("Total interfaces: %d\n", interfaces.Size())
}
