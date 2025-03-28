package uptime

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenUptime(t *testing.T) {
	ctx := sqlctx.NewContext()
	uptime, err := GenUptime(ctx)
	if err != nil {
		t.Fatalf("Failed to get uptime: %v", err)
	}

	jsonData, err := json.MarshalIndent(uptime, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal uptime: %v", err)
	}

	fmt.Printf("Uptime Results:\n%s\n", string(jsonData))
}
