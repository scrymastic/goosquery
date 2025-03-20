package listening_ports

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenListeningPorts(t *testing.T) {
	// Create context with all columns used
	ctx := context.Context{}
	// Add all possible columns to ensure they're all included in test
	ctx.Columns = []string{"pid", "port", "protocol", "family", "address", "fd", "socket", "path"}

	ports, err := GenListeningPorts(ctx)
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

func TestPort22(t *testing.T) {
	ctx := context.Context{}
	ctx.Columns = []string{"pid", "port", "protocol", "family", "address", "fd", "socket", "path"}
	ctx.AddConstant("port", "22")

	ports, err := GenListeningPorts(ctx)
	if err != nil {
		t.Fatalf("Failed to get listening ports: %v", err)
	}

	jsonData, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal listening ports to JSON: %v", err)
	}
	fmt.Printf("Listening Ports Results:\n%s\n", string(jsonData))
	fmt.Printf("Total listening ports: %d\n", len(ports))
}
