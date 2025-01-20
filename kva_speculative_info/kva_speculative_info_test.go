package kva_speculative_info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetKVASpeculativeInfo(t *testing.T) {
	info, err := GenKVASpeculativeInfo()
	if err != nil {
		t.Fatalf("Failed to get KVA speculative info: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal KVA speculative info to JSON: %v", err)
	}
	fmt.Printf("KVA Speculative Info Results:\n%s\n", string(jsonData))
}
