package programs

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenPrograms(t *testing.T) {
	programs, err := GenPrograms(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get programs: %v", err)
	}

	// Format and print JSON output
	jsonData, err := json.MarshalIndent(programs, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal programs: %v", err)
	}

	fmt.Printf("Programs Results:\n%s\n", string(jsonData))
	fmt.Printf("Total programs: %d\n", programs.Size())
}
