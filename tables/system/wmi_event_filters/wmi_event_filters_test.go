package wmi_event_filters

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWMIEventFilters(t *testing.T) {
	ctx := sqlctx.NewContext()
	filters, err := GenWMIEventFilters(ctx)
	if err != nil {
		t.Fatalf("Failed to get WMI event filters: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(filters, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal WMI event filters to JSON: %v", err)
	}
	fmt.Printf("WMI Event Filters:\n%s\n", string(jsonData))
	fmt.Printf("Total filters: %d\n", filters.Size())
}
