package services

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenServices(t *testing.T) {
	services, err := GenServices()
	if err != nil {
		t.Fatalf("Failed to get services: %v", err)
	}

	jsonData, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal services: %v", err)
	}

	t.Logf("Services JSON:\n%s", string(jsonData))
	t.Logf("Number of services: %d", len(services))
}

func ExampleGenServices() {
	services, err := GenServices()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total services: %d\n", len(services))

	jsonData, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(jsonData))
}
