package patches

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenPatches(t *testing.T) {
	patches, err := GenPatches()
	if err != nil {
		t.Fatalf("Failed to get patches: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(patches, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal patches to JSON: %v", err)
	}
	fmt.Printf("Patches Results:\n%s\n", string(jsonData))
	fmt.Printf("Total patches: %d\n", len(patches))
}
