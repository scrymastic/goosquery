package pipes

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenPipes(t *testing.T) {
	pipes, err := GenPipes()
	if err != nil {
		t.Fatalf("Failed to get pipes: %v", err)
	}

	jsonData, err := json.MarshalIndent(pipes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal pipes: %v", err)
	}

	fmt.Printf("Pipes Results:\n%s\n", string(jsonData))
	fmt.Printf("Total pipes: %d\n", len(pipes))
}
