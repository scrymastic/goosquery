package uptime

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenUptime(t *testing.T) {
	uptime, err := GenUptime(context.Context{})
	if err != nil {
		t.Fatalf("Failed to get uptime: %v", err)
	}

	jsonData, err := json.MarshalIndent(uptime, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal uptime: %v", err)
	}

	fmt.Printf("Uptime Results:\n%s\n", string(jsonData))
}
