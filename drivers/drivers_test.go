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

func TestGetDeviceList(t *testing.T) {
	devList, err := getDeviceList()
	if err != nil {
		t.Fatalf("Failed to get device list: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(devList, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal device list to JSON: %v", err)
	}
	fmt.Printf("Device List Results:\n%s\n", string(jsonData))
	fmt.Printf("Total devices: %d\n", len(devList))
}
