package process_open_sockets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenProcessOpenSockets(t *testing.T) {
	sockets, err := GenProcessOpenSockets()
	if err != nil {
		t.Fatalf("Failed to get process open sockets: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(sockets, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal open sockets: %v", err)
	}
	fmt.Printf("Process Open Sockets:\n%s\n", string(jsonData))
	fmt.Printf("Total open sockets: %d\n", len(sockets))
}
