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

	t.Logf("Number of pipes: %d", len(pipes))
	t.Logf("Named Pipes JSON:\n%s", string(jsonData))
}

func ExampleGenPipes() {
	pipes, err := GenPipes()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total pipes: %d\n", len(pipes))

	jsonData, err := json.MarshalIndent(pipes, "", "  ")
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(jsonData))
}
