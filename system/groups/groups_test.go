package groups

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenGroups(t *testing.T) {
	groups, err := GenGroups()
	if err != nil {
		t.Fatalf("Failed to get groups: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal groups to JSON: %v", err)
	}
	fmt.Printf("Groups Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", len(groups))

	// Verify we have at least some groups
	if len(groups) == 0 {
		t.Error("No groups were found, expected at least some built-in Windows groups")
	}

	// Check for some common Windows groups
	foundAdmins := false
	foundUsers := false
	for _, group := range groups {
		if group.Groupname == "Administrators" {
			foundAdmins = true
		}
		if group.Groupname == "Users" {
			foundUsers = true
		}
	}

	if !foundAdmins {
		t.Error("Did not find the Administrators group")
	}
	if !foundUsers {
		t.Error("Did not find the Users group")
	}
}
