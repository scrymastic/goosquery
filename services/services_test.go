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

	fmt.Printf("Services Results:\n%s\n", string(jsonData))
	fmt.Printf("Total services: %d\n", len(services))
}
