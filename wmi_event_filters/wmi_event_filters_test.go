package wmi_event_filters

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenWMIEventFilters(t *testing.T) {
	filters, err := GenWMIEventFilters()
	if err != nil {
		t.Fatalf("Failed to get WMI event filters: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(filters, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal WMI event filters to JSON: %v", err)
	}
	fmt.Printf("WMI Event Filters:\n%s\n", string(jsonData))
	fmt.Printf("Total filters: %d\n", len(filters))
}

func ExampleGenWMIEventFilters() {
	filters, err := GenWMIEventFilters()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Total WMI event filters: %d\n", len(filters))
}
