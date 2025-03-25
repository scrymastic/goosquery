package registry

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenRegistry(t *testing.T) {
	// Define a test registry key path
	keyPath := `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\`
	ctx := sqlctx.NewContext()
	ctx.AddConstant("key", keyPath)
	// Call the GenRegistry function
	entries, err := GenRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to get registry entries: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal registry entries to JSON: %v", err)
	}
	fmt.Printf("Registry Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", entries.Size())
}
