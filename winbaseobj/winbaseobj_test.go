package winbaseobj

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenWinBaseObj(t *testing.T) {
	objects, err := GenWinBaseObj()
	if err != nil {
		t.Fatalf("Failed to get Windows base objects: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Windows base objects to JSON: %v", err)
	}
	fmt.Printf("Windows Base Objects:\n%s\n", string(jsonData))
	fmt.Printf("Total objects: %d\n", len(objects))
}
