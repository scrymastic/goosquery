package connectivity

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenConnectivity(t *testing.T) {
	connectivity, err := GenConnectivity()
	if err != nil {
		t.Fatalf("Failed to get connectivity status: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(connectivity, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal connectivity data to JSON: %v", err)
	}
	fmt.Printf("Connectivity Status:\n%s\n", string(jsonData))

	// Basic validation
	if len(connectivity) != 1 {
		t.Errorf("Expected 1 connectivity result, got %d", len(connectivity))
	}

	// Validate that at least one connectivity state is true
	// (machine should either be connected somehow or disconnected)
	result := connectivity[0]
	hasValidState := result.Disconnected ||
		result.IPv4NoTraffic || result.IPv6NoTraffic ||
		result.IPv4Subnet || result.IPv4LocalNetwork || result.IPv4Internet ||
		result.IPv6Subnet || result.IPv6LocalNetwork || result.IPv6Internet

	if !hasValidState {
		t.Error("No valid connectivity state detected")
	}
}
