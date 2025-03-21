package appcompat_shims

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenAppCompatShims(t *testing.T) {
	shims, err := GenAppCompatShims(context.Context{})
	if err != nil {
		t.Fatalf("Failed to get AppCompat shims: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(shims, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal AppCompat shims to JSON: %v", err)
	}
	fmt.Printf("AppCompat Shims Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", len(shims))
}
