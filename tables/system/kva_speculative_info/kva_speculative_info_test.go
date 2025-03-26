package kva_speculative_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGetKVASpeculativeInfo(t *testing.T) {
	ctx := sqlctx.NewContext()
	info, err := GenKvaSpeculativeInfo(ctx)
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
