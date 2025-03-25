package windows_firewall_rules

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWindowsFirewallRules(t *testing.T) {
	// Create context with all columns used
	ctx := sqlctx.NewContext()
	// Add all possible columns to ensure they're all included in test
	ctx.Columns = []string{
		"name", "app_name", "action", "enabled", "grouping",
		"direction", "protocol", "local_addresses", "remote_addresses",
		"local_ports", "remote_ports", "icmp_types_codes",
		"profile_domain", "profile_private", "profile_public", "service_name",
	}

	rules, err := GenWindowsFirewallRules(ctx)
	if err != nil {
		t.Fatalf("Failed to get Windows Firewall Rules: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Firewall Rules to JSON: %v", err)
	}
	fmt.Printf("Windows Firewall Rules Results:\n%s\n", string(jsonData))
	fmt.Printf("Total rules: %d\n", rules.Size())

	// Basic validation of the returned data
	for i, rule := range *rules {
		if name, ok := rule.Get("name").(string); !ok || name == "" {
			t.Errorf("Rule %d: Name is empty or not a string", i)
		}
		if direction, ok := rule.Get("direction").(string); !ok || (direction != "In" && direction != "Out") {
			t.Errorf("Rule %d: Invalid direction '%v', expected 'In' or 'Out'", i, rule.Get("direction"))
		}
		if action, ok := rule.Get("action").(string); !ok || (action != "Allow" && action != "Block") {
			t.Errorf("Rule %d: Invalid action '%v', expected 'Allow' or 'Block'", i, rule.Get("action"))
		}
	}
}
