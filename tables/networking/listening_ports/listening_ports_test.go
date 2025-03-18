package listening_ports

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenListeningPorts(t *testing.T) {
	ports, err := GenListeningPorts()
	if err != nil {
		t.Fatalf("Failed to get listening ports: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal listening ports to JSON: %v", err)
	}
	fmt.Printf("Listening Ports Results:\n%s\n", string(jsonData))
	fmt.Printf("Total listening ports: %d\n", len(ports))
}
