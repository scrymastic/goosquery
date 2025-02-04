package windows_firewall_rules

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenWindowsFirewallRules(t *testing.T) {
	rules, err := GenWindowsFirewallRules()
	if err != nil {
		t.Fatalf("Failed to get Windows Firewall Rules: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Firewall Rules to JSON: %v", err)
	}
	fmt.Printf("Windows Firewall Rules Results:\n%s\n", string(jsonData))
	fmt.Printf("Total rules: %d\n", len(rules))

	// Basic validation of the returned data
	for i, rule := range rules {
		if rule.Name == "" {
			t.Errorf("Rule %d: Name is empty", i)
		}
		if rule.Direction != "In" && rule.Direction != "Out" {
			t.Errorf("Rule %d: Invalid direction '%s', expected 'In' or 'Out'", i, rule.Direction)
		}
		if rule.Action != "Allow" && rule.Action != "Block" {
			t.Errorf("Rule %d: Invalid action '%s', expected 'Allow' or 'Block'", i, rule.Action)
		}
	}
}
