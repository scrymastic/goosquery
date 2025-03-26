package windows_update_history

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWindowsUpdateHistory(t *testing.T) {
	ctx := sqlctx.NewContext()
	history, err := GenWindowsUpdateHistory(ctx)
	if err != nil {
		t.Fatalf("Failed to get Windows update history: %v", err)
	}

	if history.Size() == 0 {
		t.Fatal("No Windows update history entries found")
	}

	// Print the first entry's date value for debugging
	jsonData, err := json.MarshalIndent(history.GetRow(0), "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Windows update history to JSON: %v", err)
	}
	fmt.Printf("Windows Update History Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", history.Size())
}
