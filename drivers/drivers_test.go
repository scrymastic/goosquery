package drivers

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenDrivers(t *testing.T) {
	drivers, err := GenDrivers()
	if err != nil {
		t.Fatalf("Failed to get drivers: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(drivers, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal drivers to JSON: %v", err)
	}
	fmt.Printf("Drivers Results:\n%s\n", string(jsonData))
	fmt.Printf("Total drivers: %d\n", len(drivers))
}

func ExampleGenDrivers() {
	drivers, err := GenDrivers()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Total drivers: %d\n", len(drivers))
}
