package hash

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenHash(t *testing.T) {
	file := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	hash, err := GenFileHash(sqlctx.NewContext(), file)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}
	// Print results as JSON for visibility
	jsonData, err := json.MarshalIndent(hash, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal hash to JSON: %v", err)
	}
	fmt.Printf("Hash Results:\n%s\n", string(jsonData))
}
