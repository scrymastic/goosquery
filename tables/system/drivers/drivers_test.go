package drivers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenDrivers(t *testing.T) {
	drivers, err := GenDrivers(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get drivers: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(drivers, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal drivers to JSON: %v", err)
	}
	fmt.Printf("Drivers Results:\n%s\n", string(jsonData))
	fmt.Printf("Total drivers: %d\n", drivers.Size())
}
