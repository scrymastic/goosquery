package logical_drives

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenLogicalDrives(t *testing.T) {
	ctx := sqlctx.NewContext()
	drives, err := GenLogicalDrives(ctx)
	if err != nil {
		t.Fatalf("Failed to get logical drives: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(drives, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal logical drives to JSON: %v", err)
	}
	fmt.Printf("Logical Drives Results:\n%s\n", string(jsonData))
	fmt.Printf("Total drives: %d\n", drives.Size())
}
